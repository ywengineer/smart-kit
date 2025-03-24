package main

import (
	"github.com/bytedance/sonic"
	"github.com/ywengineer/smart-kit/passport/internal/model"
	"github.com/ywengineer/smart/utility"
	"gopkg.in/yaml.v3"
	"net"
	"testing"
)

func TestYamlConf(t *testing.T) {
	ym, _ := yaml.Marshal(Configuration{
		RDB:  utility.RdbProperties{},
		Cors: &Cors{},
	})
	t.Log(string(ym))
	//
	t.Log(sonic.MarshalString(model.PassportBinding{}))
}

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
