package hw

import (
	"context"
	"fmt"
	"net/url"
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
	config     Config
	token      *tokenResponse
	mu         sync.RWMutex
	lastUpdate int64 // 令牌最后更新时间（秒级时间戳）
}

// tokenResponse 令牌响应结构体（对应官方文档）
type tokenResponse struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
	Ttl         int    `json:"ttl"`
	Timestamp   string `json:"timestamp"`
}

// NewTokenManager 初始化令牌管理器
func NewTokenManager(config Config) (*TokenManager, error) {
	tm := &TokenManager{config: config}
	if err := tm.refreshToken(); err != nil {
		return nil, errors.WithMessage(err, "failed to init huawei token manager")
	}
	return tm, nil
}

// getToken 安全获取令牌（自动判断是否需要刷新）
func (tm *TokenManager) getToken() (string, error) {
	tm.mu.RLock()
	// 检查令牌是否有效（未过期且距离过期>缓冲时间）
	isValid := tm.token != nil && time.Now().Unix()-tm.lastUpdate < int64(tm.token.Ttl-tokenExpiryBuffer)
	tm.mu.RUnlock()

	if !isValid {
		//
		if err := tm.refreshToken(); err != nil {
			return "", err
		}
	}
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.token.AccessToken, nil
}

// refreshToken 刷新令牌（核心：签名生成+请求发送）
func (tm *TokenManager) refreshToken() error {
	tokenURL := prodTokenURL
	ts := time.Now()
	urlValue := url.Values{"grant_type": {"client_credentials"}, "client_secret": {tm.config.ClientSecret}, "client_id": {tm.config.ClientID}}
	//
	statusCode, resp, err := rpcs.GetDefaultRpc().Post(context.Background(), rpcs.ContentTypeUrlencoded, tokenURL, nil, rpcs.BytesBody{V: []byte(urlValue.Encode())})
	if err != nil {
		return errors.WithMessage(err, "Failed to get huawei token")
	} else if statusCode != consts.StatusOK {
		return errors.New(fmt.Sprintf("Failed to get huawei token with status: %d, body = %s", statusCode, string(resp)))
	}
	var tokenResp tokenResponse
	if err := sonic.Unmarshal(resp, &tokenResp); err != nil || !strings.EqualFold(tokenResp.Code, "OK") {
		return errors.WithMessagef(err, "Failed to parse the huawei token response: %s", string(resp))
	}
	//
	tm.mu.Lock()
	defer tm.mu.Unlock()
	//
	tokenResp.AccessToken = encodeAccessToken(tokenResp.AccessToken)
	//
	tm.token, tm.lastUpdate = &tokenResp, ts.Unix()
	return nil
}
