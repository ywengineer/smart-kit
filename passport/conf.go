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
