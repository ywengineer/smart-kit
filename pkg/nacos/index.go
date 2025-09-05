package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
)

type Nacos struct {
	Ip          string `json:"ip" yaml:"ip" env:"NACOS_IP"`
	Port        uint64 `json:"port" yaml:"port" env:"NACOS_PORT"`
	ContextPath string `json:"context_path" yaml:"context-path" env:"NACOS_CONTEXT_PATH"`
	TimeoutMs   uint64 `json:"timeout_ms" yaml:"timeout-ms" env:"NACOS_TIMEOUT_MS"`
	Namespace   string `json:"namespace" yaml:"namespace" env:"NACOS_NAMESPACE"`
	User        string `json:"user" yaml:"user" env:"NACOS_USER"`
	Password    string `json:"password" yaml:"password" env:"NACOS_PASSWORD"`
	Cluster     string `json:"cluster" yaml:"cluster" env:"NACOS_CLUSTER"`
	Group       string `json:"group" yaml:"group" env:"NACOS_GROUP"`
}

// NewNamingClientWithConfig
// the logLevel must be one of debug,info,warn,error, default value is debug
func NewNamingClientWithConfig(c Nacos, logLevel string) (naming_client.INamingClient, error) {
	return NewNacosNamingClient(c.Ip, c.Port, c.ContextPath, c.TimeoutMs, c.Namespace, c.User, c.Password, logLevel)
}

// NewConfigClientWithConfig
// the logLevel must be one of debug,info,warn,error, default value is debug
func NewConfigClientWithConfig(c Nacos, logLevel string) (config_client.IConfigClient, error) {
	return NewNacosConfigClient(c.Ip, c.Port, c.ContextPath, c.TimeoutMs, c.Namespace, c.User, c.Password, logLevel)
}
