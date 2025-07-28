package apps

import (
	"fmt"
	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"gitee.com/ywengineer/smart-kit/pkg/nacos"
	"gitee.com/ywengineer/smart-kit/pkg/oauths"
	"gitee.com/ywengineer/smart-kit/pkg/rdbs"
	"gitee.com/ywengineer/smart-kit/pkg/rpcs"
	"time"
)

type Configuration struct {
	Port             int                `json:"port" yaml:"port"`
	BasePath         string             `json:"base_path" yaml:"base-path"`
	RDB              rdbs.Properties    `json:"rdb" yaml:"rdb"`
	Redis            string             `yaml:"redis" json:"redis"` // redis://user:password@host:port/?db=0&node=host:port&node=host:port
	DistributeLock   bool               `json:"distribute_lock" yaml:"distribute-lock"`
	MaxRequestBodyKB int                `json:"max_request_body_kb,omitempty" yaml:"max-request-body-kb,omitempty"`
	Cors             *Cors              `json:"cors,omitempty" yaml:"cors,omitempty"`
	Jwt              *JwtConfig         `json:"jwt,omitempty" yaml:"jwt,omitempty"`
	LogLevel         logk.Level         `json:"log_level" yaml:"log-level"`
	TraceLevel       TraceLevel         `json:"trace_level" yaml:"trace-level"`
	Nacos            *nacos.Nacos       `json:"nacos,omitempty" yaml:"nacos,omitempty"`
	RegistryInfo     *ServiceInfo       `json:"registry_info" yaml:"registry-info"`
	DiscoveryEnable  bool               `json:"discovery_enable" yaml:"discovery-enable"`
	RpcClientInfo    rpcs.RpcClientInfo `json:"rpc_client_info" yaml:"rpc-client-info"`
	OAuth            oauths.Oauth       `json:"oauth" yaml:"oauth"`
	SignKey          string             `json:"sign_key" yaml:"sign-key"`
	Profile          Profiling          `json:"profile" yaml:"profile"`
	AccessLog        string             `json:"access_log" yaml:"access-log"`
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

func (s *ServiceInfo) String() string {
	return fmt.Sprintf("%s(%s)[%d]", s.ServiceName, s.Addr, s.Weight)
}
