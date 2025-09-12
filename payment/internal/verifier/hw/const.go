package hw

import (
	"encoding/base64"
	"fmt"
)

const (
	// 令牌接口地址
	prodTokenURL = "https://oauth-login.cloud.huawei.com/oauth2/v3/token"
	// 支付检验接口地址（按 invoiceId 查询）
	prodVerifyURL    = "https://public-api.rustore.ru/public/v2/purchase?invoiceId="
	sandboxVerifyURL = "https://public-api.rustore.ru/public/sandbox/v2/purchase?invoiceId="
	// 固定参数
	tokenExpiryBuffer = 30 // 令牌刷新缓冲时间（提前30秒刷新，避免过期）
)

// Config config
type Config struct {
	ClientID     string // 控制台获取的 Client ID
	ClientSecret string // 控制台获取的 application public key, base64 encode
	IsSandbox    bool   // 是否启用沙箱环境
}

func encodeAccessToken(t string) string {
	oriString := fmt.Sprintf("APPAT:%s", t)
	var authString = base64.StdEncoding.EncodeToString([]byte(oriString))
	var authHeaderString = fmt.Sprintf("Basic %s", authString)
	return authHeaderString
}
