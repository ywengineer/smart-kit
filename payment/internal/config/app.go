package config

import (
	"context"

	"gitee.com/ywengineer/smart-kit/pkg/loaders"
	"gitee.com/ywengineer/smart-kit/pkg/nacos"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var p Payment
var loader loaders.SmartLoader

type Payment struct {
	Auth Auth `json:"auth" yaml:"auth" redis:"auth"`
}

func Watch(ctx context.Context, n nacos.Nacos) error {
	nc, err := nacos.NewConfigClientWithConfig(n, "info")
	if err != nil {
		return err
	}
	loader = loaders.NewNacosLoader(nc, "", "payment.yaml", loaders.NewYamlDecoder())
	err = loader.Load(&p)
	if err != nil {
		return err
	}
	return loader.Watch(ctx, func(data string) error {
		err := loader.Unmarshal([]byte(data), &p)
		hlog.CtxInfof(ctx, "payment application config change: %+v", p)
		return err
	})
}

func Get() Payment {
	return p
}
