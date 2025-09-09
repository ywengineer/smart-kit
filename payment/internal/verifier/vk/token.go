package vk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RustoreConfig rustore config
type RustoreConfig struct {
	ClientID     string // 控制台获取的 Client ID
	ClientSecret string // 控制台获取的 Client Secret
	IsSandbox    bool   // 是否启用沙箱环境
}

// TokenManager token manager
type TokenManager struct {
	config     RustoreConfig
	token      *TokenResponse
	mu         sync.RWMutex
	lastUpdate int64 // 令牌最后更新时间（秒级时间戳）
}

// TokenResponse 令牌响应结构体（对应官方文档）
type TokenResponse struct {
	AccessToken string `json:"access_token"` // JWE 令牌
	TokenType   string `json:"token_type"`   // 固定为 "bearer"
	ExpiresIn   int    `json:"expires_in"`   // 有效期（秒）
}

// NewTokenManager 初始化令牌管理器
func NewTokenManager(config RustoreConfig) (*TokenManager, error) {
	tm := &TokenManager{config: config}
	if err := tm.refreshToken(); err != nil {
		return nil, fmt.Errorf("failed to init token manager：%w", err)
	}
	return tm, nil
}

// getToken 安全获取令牌（自动判断是否需要刷新）
func (tm *TokenManager) getToken() (string, error) {
	tm.mu.RLock()
	// 检查令牌是否有效（未过期且距离过期>缓冲时间）
	isValid := tm.token != nil && time.Now().Unix()-tm.lastUpdate < int64(tm.token.ExpiresIn-tokenExpiryBuffer)
	tm.mu.RUnlock()

	if !isValid {
		// 令牌无效，触发刷新（写锁）
		if err := tm.refreshToken(); err != nil {
			return "", fmt.Errorf("failed to refresh token：%w", err)
		}
	}
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.token.AccessToken, nil
}

// refreshToken 刷新令牌（核心：签名生成+请求发送）
func (tm *TokenManager) refreshToken() error {
	tokenURL := prodTokenURL
	if tm.config.IsSandbox {
		tokenURL = sandboxTokenURL
	}
	timestamp := time.Now().UnixMilli()
	signStr := fmt.Sprintf("client_id=%s&timestamp=%d", tm.config.ClientID, timestamp)
	h := hmac.New(sha256.New, []byte(tm.config.ClientSecret))
	if _, err := h.Write([]byte(signStr)); err != nil {
		return fmt.Errorf("签名计算失败：%w", err)
	}
	signature := hex.EncodeToString(h.Sum(nil))
	// 4. 构造请求体
	reqBody := struct {
		GrantType     string `json:"grant_type"`
		ClientID      string `json:"client_id"`
		Timestamp     int64  `json:"timestamp"`
		Signature     string `json:"signature"`
		SignAlgorithm string `json:"sign_algorithm"`
	}{
		GrantType:     grantType,
		ClientID:      tm.config.ClientID,
		Timestamp:     timestamp,
		Signature:     signature,
		SignAlgorithm: signAlgorithm,
	}
	// 5. 发送 POST 请求
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return fmt.Errorf("请求创建失败：%w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(reqBodyBytes)))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求发送失败：%w", err)
	}
	defer resp.Body.Close()

	// 6. 解析响应（区分成功/失败）
	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("错误响应解析失败：%w", err)
		}
		return fmt.Errorf("令牌接口返回错误：%s（%s）", errResp.Error, errResp.ErrorDescription)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("令牌响应解析失败：%w", err)
	}
	if tokenResp.AccessToken == "" {
		return fmt.Errorf("获取到空令牌")
	}

	// 7. 更新令牌（写锁保护）
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.token = &tokenResp
	tm.lastUpdate = time.Now().Unix()

	return nil
}
