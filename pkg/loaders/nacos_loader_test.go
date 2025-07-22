package loaders

import (
	"context"
	"gitee.com/ywengineer/smart-kit/pkg/nacos"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNacosLoader(t *testing.T) {
	nc, err := nacos.NewNacosConfigClient("192.168.44.128", 8848, "/nacos", 5000,
		"a7aabc24-17a7-4ac5-978f-6f933ce19dd4", "nacos", "nacos", "debug")
	assert.Nil(t, err)
	//
	c := &Conf{}
	loader := NewNacosLoader(nc, "DEFAULT_GROUP", "smart.gate.yaml", NewYamlDecoder())
	err = loader.Load(c)
	assert.Nil(t, err)
	t.Logf("%v", *c)
	err = loader.Watch(context.Background(), func(conf string) error {
		_ = loader.Unmarshal([]byte(conf), c)
		t.Logf("config change: %v", *c)
		return nil
	})
	assert.Nil(t, err)
	<-utilk.WatchQuitSignal()
	t.Log("test finished")
}
