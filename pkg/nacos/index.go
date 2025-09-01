package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
)

type Nacos struct {
	Ip          string `json:"ip" yaml:"ip"`
	Port        uint64 `json:"port" yaml:"port"`
	ContextPath string `json:"context_path" yaml:"context-path"`
	TimeoutMs   uint64 `json:"timeout_ms" yaml:"timeout-ms"`
	Namespace   string `json:"namespace" yaml:"namespace"`
	User        string `json:"user" yaml:"user"`
	Password    string `json:"password" yaml:"password"`
	Cluster     string `json:"cluster" yaml:"cluster"`
	Group       string `json:"group" yaml:"group"`
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
