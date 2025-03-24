package main

import (
	"github.com/ywengineer/smart-kit/passport/pkg/middleware"
	"github.com/ywengineer/smart/utility"
	"go.uber.org/zap/zapcore"
	"time"
)

type Configuration struct {
	Port             int                   `json:"port" yaml:"port"`
	BasePath         string                `json:"base_path" yaml:"base-path"`
	RDB              utility.RdbProperties `json:"rdb" yaml:"rdb"`
	Redis            string                `yaml:"redis" json:"redis"` // redis://user:password@host:port/?db=0&node=host:port&node=host:port
	DistributeLock   bool                  `json:"distribute_lock" yaml:"distribute-lock"`
	MaxRequestBodyKB int                   `json:"max_request_body_kb,omitempty" yaml:"max-request-body-kb,omitempty"`
	Cors             *Cors                 `json:"cors,omitempty" yaml:"cors,omitempty"`
	Jwt              *middleware.JwtConfig `json:"jwt,omitempty" yaml:"jwt,omitempty"`
	LogLevel         zapcore.Level         `json:"log_level" yaml:"log-level"`
	Nacos            *Nacos                `json:"nacos,omitempty" yaml:"nacos,omitempty"`
	RegistryInfo     *ServiceInfo          `json:"registry_info" yaml:"registry-info"`
	DiscoveryEnable  bool                  `json:"discovery_enable" yaml:"discovery-enable"`
}

type Cors struct {
	AllowOrigins     []string      `json:"allow_origins" yaml:"allow-origins"`
	AllowMethods     []string      `json:"allow_methods" yaml:"allow-methods"`
	AllowHeaders     []string      `json:"allow_headers" yaml:"allow-headers"`
	AllowCredentials bool          `json:"allow_credentials" yaml:"allow-credentials"`
	ExposeHeaders    []string      `json:"expose_headers" yaml:"expose-headers"`
	MaxAge           time.Duration `json:"max_age" yaml:"max-age"`
	AllowWildcard    bool          `json:"allow_wildcard" yaml:"allow-wildcard"`
}

type Nacos struct {
	Ip          string `json:"ip" yaml:"ip"`
	Port        uint64 `json:"port" yaml:"port"`
	ContextPath string `json:"context_path" yaml:"context-path"`
	TimeoutMs   uint64 `json:"timeout_ms" yaml:"timeout-ms"`
	Namespace   string `json:"namespace" yaml:"namespace"`
	User        string `json:"user" yaml:"user"`
	Password    string `json:"password" yaml:"password"`
}
type ServiceInfo struct {
	// ServiceName will be set in hertz by default
	ServiceName string `json:"service_name" yaml:"service-name"`
	// Addr will be set in hertz by default
	Addr string `json:"addr" yaml:"addr"`
	// Weight will be set in hertz by default
	Weight int `json:"weight" yaml:"weight"`
	// extend other infos with Tags.
	Tags map[string]string `json:"tags" yaml:"tags"`
}
