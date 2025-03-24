package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// NewNacosConfigClient
// contextPath, nacos server context path
// the logLevel must be debug,info,warn,error, default value is info
func NewNacosConfigClient(ipAddr string, port uint64, contextPath string,
	timeoutMs uint64,
	namespace, user, password, logLevel string,
) (config_client.IConfigClient, error) {
	// create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ipAddr, port, constant.WithContextPath(contextPath)),
	}
	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(namespace),
		constant.WithTimeoutMs(timeoutMs),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithUsername(user),
		constant.WithPassword(password),
		constant.WithLogLevel(logLevel),
	)
	// create loader client
	return clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
}

// NewNacosNamingClient
// contextPath, nacos server context path
// the logLevel must be debug,info,warn,error, default value is info
func NewNacosNamingClient(ipAddr string, port uint64, contextPath string,
	timeoutMs uint64,
	namespace, user, password, logLevel string,
) (naming_client.INamingClient, error) {
	// create ServerConfig
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ipAddr, port, constant.WithContextPath(contextPath)),
	}
	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(namespace),
		constant.WithTimeoutMs(timeoutMs),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithUsername(user),
		constant.WithPassword(password),
		constant.WithLogLevel(logLevel),
	)
	// create loader client
	return clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
}
