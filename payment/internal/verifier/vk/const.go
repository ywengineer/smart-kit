package vk

const (
	// 令牌接口地址
	prodTokenURL    = "https://api.rustore.ru/auth/jwe"
	sandboxTokenURL = "https://api.rustore.ru/sandbox/auth/jwe"
	// 支付检验接口地址（按 invoiceId 查询）
	prodCheckURL    = "https://public-api.rustore.ru/public/v2/purchase/%d"
	sandboxCheckURL = "https://public-api.rustore.ru/public/sandbox/v2/purchase/%d"
	// 固定参数
	grantType         = "client_credentials"
	signAlgorithm     = "HMAC-SHA256"
	tokenExpiryBuffer = 30 // 令牌刷新缓冲时间（提前30秒刷新，避免过期）
)
