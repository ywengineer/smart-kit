package nets

import (
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"net"
)

const UnknownIpAddr = "-"

func GetDefaultIpv4() string {
	return GetDefaultIp(true)
}

func GetDefaultIp(v4 bool) string {
	inters, err := net.Interfaces()
	if err != nil {
		return UnknownIpAddr
	}
	for _, inter := range inters {
		if inter.Flags&net.FlagLoopback != net.FlagLoopback &&
			inter.Flags&net.FlagUp != 0 {
			addrs, err := inter.Addrs()
			if err != nil {
				return UnknownIpAddr
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if v4 && ipnet.IP.To4() == nil {
						continue
					}
					return ipnet.IP.String()
				}
			}
		}
	}
	return UnknownIpAddr
}

// Is2xx 用于检查状态码是否为 2xx
func Is2xx(statusCode int) bool {
	return statusCode >= consts.StatusOK && statusCode < consts.StatusMultipleChoices
}
