package vk

import "github.com/samber/lo"

const (
	// 令牌接口地址
	prodTokenURL = "https://public-api.rustore.ru/public/auth/"
	// 支付检验接口地址（按 invoiceId 查询）
	prodVerifyURL    = "https://public-api.rustore.ru/public/v2/purchase/"
	sandboxVerifyURL = "https://public-api.rustore.ru/public/sandbox/v2/purchase/"

	// 固定参数
	tokenExpiryBuffer = 30 // 令牌刷新缓冲时间（提前30秒刷新，避免过期）
)

// RustoreConfig rustore config
type RustoreConfig struct {
	ClientID     string // 控制台获取的 Client ID
	ClientSecret string // 控制台获取的 Client Secret
	IsSandbox    bool   // 是否启用沙箱环境
	Apps         []string
}

func (rc RustoreConfig) IsValidApp(appId string) bool {
	return len(rc.Apps) > 0 && lo.Contains(rc.Apps, appId)
}
