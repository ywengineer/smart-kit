package apps

import (
	"fmt"
	"time"

	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"gitee.com/ywengineer/smart-kit/pkg/nacos"
	"gitee.com/ywengineer/smart-kit/pkg/oauths"
	"gitee.com/ywengineer/smart-kit/pkg/rdbs"
	"gitee.com/ywengineer/smart-kit/pkg/rpcs"
)

type Configuration struct {
	Port             int                `json:"port" yaml:"port" mapstructure:"port" env:"SERVICE_PORT"`
	BasePath         string             `json:"base_path" yaml:"base-path" mapstructure:"base-path" env:"SERVICE_BASE_PATH"`
	RDB              rdbs.Properties    `json:"rdb" yaml:"rdb" mapstructure:"rdb"`
	Redis            string             `yaml:"redis" json:"redis" mapstructure:"redis" env:"SERVICE_REDIS"` // redis://user:password@host:port/?db=0&node=host:port&node=host:port
	DistributeLock   bool               `json:"distribute_lock" yaml:"distribute-lock" mapstructure:"distribute-lock"`
	MaxRequestBodyKB int                `json:"max_request_body_kb,omitempty" yaml:"max-request-body-kb,omitempty" mapstructure:"max-request-body-kb,omitempty"`
	Cors             *Cors              `json:"cors,omitempty" yaml:"cors,omitempty" mapstructure:"cors,omitempty"`
	Jwt              *JwtConfig         `json:"jwt,omitempty" yaml:"jwt,omitempty" mapstructure:"jwt,omitempty"`
	LogLevel         logk.Level         `json:"log_level" yaml:"log-level" mapstructure:"log-level"`
	TraceLevel       TraceLevel         `json:"trace_level" yaml:"trace-level" mapstructure:"trace-level"`
	Nacos            *nacos.Nacos       `json:"nacos,omitempty" yaml:"nacos,omitempty" mapstructure:"nacos,omitempty"`
	RegistryInfo     *ServiceInfo       `json:"registry_info,omitempty" yaml:"registry-info,omitempty" mapstructure:"registry-info,omitempty"`
	DiscoveryEnable  bool               `json:"discovery_enable" yaml:"discovery-enable" mapstructure:"discovery-enable"`
	RpcClientInfo    rpcs.RpcClientInfo `json:"rpc_client_info" yaml:"rpc-client-info" mapstructure:"rpc-client-info"`
	OAuth            oauths.Oauth       `json:"oauth" yaml:"oauth" mapstructure:"oauth"`
	SignKey          string             `json:"sign_key" yaml:"sign-key" mapstructure:"sign-key"`
	Profile          Profiling          `json:"profile" yaml:"profile" mapstructure:"profile"`
	AccessLog        string             `json:"access_log" yaml:"access-log" mapstructure:"access-log"`
	RateLimitEnabled bool               `json:"rate_limit_enabled" yaml:"rate-limit-enabled" mapstructure:"rate-limit-enabled"`
}

type Cors struct {
	AllowOrigins     []string      `json:"allow_origins" yaml:"allow-origins" mapstructure:"allow-origins"`
	AllowMethods     []string      `json:"allow_methods" yaml:"allow-methods" mapstructure:"allow-methods"`
	AllowHeaders     []string      `json:"allow_headers" yaml:"allow-headers" mapstructure:"allow-headers"`
	AllowCredentials bool          `json:"allow_credentials" yaml:"allow-credentials" mapstructure:"allow-credentials"`
	ExposeHeaders    []string      `json:"expose_headers" yaml:"expose-headers" mapstructure:"expose-headers"`
	MaxAge           time.Duration `json:"max_age" yaml:"max-age" mapstructure:"max-age"`
	AllowWildcard    bool          `json:"allow_wildcard" yaml:"allow-wildcard" mapstructure:"allow-wildcard"`
}

type ServiceInfo struct {
	// ServiceName will be set in hertz by default
	ServiceName string `json:"service_name" yaml:"service-name" mapstructure:"service-name" env:"SERVICE_NAME"`
	// Addr will be set in hertz by default
	Addr string `json:"addr" yaml:"addr" mapstructure:"addr" env:"SERVICE_ADDR"`
	// Weight will be set in hertz by default
	Weight int `json:"weight" yaml:"weight" mapstructure:"weight" env:"SERVICE_WEIGHT"`
	// extend other infos with Tags.
	Tags map[string]string `json:"tags" yaml:"tags" mapstructure:"tags" env:"SERVICE_TAGS"`
}

func (s *ServiceInfo) String() string {
	return fmt.Sprintf("%s(%s)[%d]", s.ServiceName, s.Addr, s.Weight)
}
