package config

type ChannelProperty struct {
	Validator    string   `json:"validator" yaml:"validator"`         // rustore, huawei, xiaomi
	ClientID     string   `json:"client_id" yaml:"client-id"`         // 控制台获取的 Client ID
	ClientSecret string   `json:"client_secret" yaml:"client-secret"` // 控制台获取的 Client Secret
	Sandbox      bool     `json:"sandbox" yaml:"sandbox"`             // 是否启用沙箱环境
	Apps         []string `json:"apps" yaml:"apps"`
}
