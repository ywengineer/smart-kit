package config

import (
	"context"

	"gitee.com/ywengineer/smart-kit/pkg/loaders"
	"gitee.com/ywengineer/smart-kit/pkg/nacos"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var mt *Metadata
var p Payment
var loader loaders.SmartLoader

func init() {
	mt = &Metadata{
		channelMap:    make(map[string]Channel),
		productMap:    make(map[uint64]Product),
		gameServerMap: make(map[uint64]GameServerInfo),
	}
}

type RemoteUrl struct {
	EnableUpdate bool   `json:"enableUpdate" yaml:"enableUpdate"`
	Product      string `yaml:"product" json:"product"`
	GameServer   string `yaml:"gameServer" json:"gameServer"`
	Platform     string `yaml:"platform" json:"platform"`
}

type Queue struct {
	Workers int            `yaml:"workers" json:"workers"`
	Queues  map[string]int `yaml:"queues" json:"queues"`
}

type Payment struct {
	Auth      Auth      `json:"auth" yaml:"auth" redis:"auth"`
	RemoteUrl RemoteUrl `json:"remoteUrl" yaml:"remoteUrl" redis:"remoteUrl"`
	Queue     Queue     `json:"queue" yaml:"queue" redis:"queue"`
}

type Listener func(c Payment)

func Watch(ctx context.Context, n nacos.Nacos, l Listener) error {
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
		hlog.CtxInfof(ctx, "payment application config change: %+v, error: %v", p, err)
		if err == nil {
			l(p)
		}
		return err
	})
}

func Get() Payment {
	return p
}

func GetMeta() *Metadata {
	return mt
}
