package apps

import (
	"net"
	"testing"
)

func TestIpAddr(t *testing.T) {
	inters, err := net.Interfaces()
	if err != nil {
		t.Fatalf("UNKNOWN_IP_ADDR")
	}
	for _, inter := range inters {
		if inter.Flags&net.FlagLoopback != net.FlagLoopback &&
			inter.Flags&net.FlagUp != 0 {
			addrs, err := inter.Addrs()
			if err != nil {
				t.Fatalf("UNKNOWN_IP_ADDR")
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					t.Logf("ip addr: %s, private = %v", ipnet.IP.String(), ipnet.IP.IsPrivate())
				}
			}
		}
	}
}
