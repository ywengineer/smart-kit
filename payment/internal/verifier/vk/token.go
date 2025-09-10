package vk

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"gitee.com/ywengineer/smart-kit/pkg/rpcs"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/pkg/errors"
)

// TokenManager token manager
type TokenManager struct {
	config     RustoreConfig
	token      *tokenResponse
	key        *rsa.PrivateKey
	mu         sync.RWMutex
	lastUpdate int64 // 令牌最后更新时间（秒级时间戳）
}

// tokenResponse 令牌响应结构体（对应官方文档）
type tokenResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Body    struct {
		Jwe string `json:"jwe"`
		Ttl int    `json:"ttl"`
	} `json:"body"`
	Timestamp string `json:"timestamp"`
}

// NewTokenManager 初始化令牌管理器
func NewTokenManager(config RustoreConfig) (*TokenManager, error) {
	tm := &TokenManager{config: config}
	if err := tm.parsePrivateKey(); err != nil {
		return nil, err
	}
	if err := tm.refreshToken(); err != nil {
		return nil, errors.WithMessage(err, "failed to init rustore token manager")
	}
	return tm, nil
}

func (tm *TokenManager) parsePrivateKey() error {
	// private key bytes
	privateKeyBytes, err := base64.StdEncoding.DecodeString(tm.config.ClientSecret)
	if err != nil {
		return errors.WithMessage(err, "failed to decode rustore private key")
	}
	//
	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return errors.WithMessage(err, "failed to parse rustore private key as x509 format")
	}
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return errors.New("Not a valid RSA private key")
	}
	tm.key = rsaPrivateKey
	return nil
}

// getToken 安全获取令牌（自动判断是否需要刷新）
func (tm *TokenManager) getToken() (string, error) {
	tm.mu.RLock()
	// 检查令牌是否有效（未过期且距离过期>缓冲时间）
	isValid := tm.token != nil && time.Now().Unix()-tm.lastUpdate < int64(tm.token.Body.Ttl-tokenExpiryBuffer)
	tm.mu.RUnlock()

	if !isValid {
		//
		if err := tm.refreshToken(); err != nil {
			return "", err
		}
	}
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.token.Body.Jwe, nil
}

func (tm *TokenManager) sign(t string) (string, error) {
	// 计算SHA-512哈希
	hash := crypto.SHA512.New()
	hash.Write([]byte(tm.config.ClientID + t))
	hashed := hash.Sum(nil)
	// 使用RSA私钥签名
	signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, tm.key, crypto.SHA512, hashed)
	if err != nil {
		return "", errors.WithMessage(err, "Signature failed")
	}
	// 对签名结果进行Base64编码
	return base64.StdEncoding.EncodeToString(signatureBytes), nil
}

// refreshToken 刷新令牌（核心：签名生成+请求发送）
func (tm *TokenManager) refreshToken() error {
	tokenURL := prodTokenURL
	ts := time.Now()
	timestamp := ts.Format(time.RFC3339)
	signature, err := tm.sign(timestamp)
	if err != nil {
		return err
	}
	//
	reqBody := struct {
		KeyId     string `json:"keyId"`
		Timestamp string `json:"timestamp"`
		Signature string `json:"signature"`
	}{
		KeyId:     tm.config.ClientID,
		Timestamp: timestamp,
		Signature: signature,
	}
	//
	statusCode, resp, err := rpcs.GetDefaultRpc().Post(context.Background(), rpcs.ContentTypeJSON, tokenURL, nil, rpcs.JsonBody{V: reqBody})
	if err != nil {
		return errors.WithMessage(err, "Failed to get rustore token")
	} else if statusCode != consts.StatusOK {
		return errors.New(fmt.Sprintf("Failed to get rustore token with status: %d", statusCode))
	}
	var tokenResp tokenResponse
	if err := sonic.Unmarshal(resp, &tokenResp); err != nil || !strings.EqualFold(tokenResp.Code, "OK") {
		return errors.WithMessagef(err, "Failed to parse the rustore token response: %s", string(resp))
	}
	//
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.token, tm.lastUpdate = &tokenResp, ts.Unix()
	return nil
}
