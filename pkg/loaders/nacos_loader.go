package loaders

import (
	"context"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/pkg/errors"
	"log"
	"time"
)

type nacosLoader struct {
	nc      config_client.IConfigClient
	group   string
	dataId  string
	decoder Decoder
}

func NewNacosLoader(nacos config_client.IConfigClient, group string, dataId string, decoder Decoder) SmartLoader {
	return &nacosLoader{
		nc:      nacos,
		dataId:  utilk.DefaultIfEmpty(dataId, "smart.server.yaml"),
		group:   utilk.DefaultIfEmpty(group, "DEFAULT_GROUP"),
		decoder: decoder,
	}
}

func NewDefaultNacosLoader(nacos config_client.IConfigClient, dataId string, decoder Decoder) SmartLoader {
	return NewNacosLoader(nacos, "DEFAULT_GROUP", dataId, decoder)
}

func (nl *nacosLoader) Unmarshal(data []byte, out interface{}) error {
	return nl.decoder.Unmarshal(data, out)
}

func (nl *nacosLoader) Load(out interface{}) error {
	if err := nl.check(); err != nil {
		return err
	}
	// get loader
	content, err := nl.nc.GetConfig(vo.ConfigParam{Group: nl.group, DataId: nl.dataId})
	if err != nil {
		return errors.WithMessage(err, "load loader content from nacos error")
	}
	return nl.Unmarshal([]byte(content), out)
}

func (nl *nacosLoader) check() error {
	if nl.nc == nil {
		return errors.New("nacos client have not been initialized.")
	}
	if nl.decoder == nil {
		return errors.New("nil loader decoder is not allowed")
	}
	if len(nl.group) == 0 || len(nl.dataId) == 0 {
		return errors.New("empty dataId and group is not allowed")
	}
	return nil
}

func (nl *nacosLoader) Watch(ctx context.Context, callback WatchCallback) error {
	if err := nl.check(); err != nil {
		return err
	}
	p := vo.ConfigParam{
		DataId: nl.dataId,
		Group:  nl.group,
		OnChange: func(namespace, group, dataId, data string) {
			_ = callback(data)
		},
	}
	go func() {
		defer func() {
			_ = nl.nc.CancelListenConfig(p)
			nl.nc.CloseClient()
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := ctx.Err(); err != nil {
					log.Printf("[nacosLoader] nacos loader watcher stopped. encounter an error: %v\n", err)
					return
				}
				time.Sleep(time.Second)
			}
		}
	}()
	return nl.nc.ListenConfig(p)
}
