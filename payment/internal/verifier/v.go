package verifier

import (
	"sync"

	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/hw"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/inf"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/vk"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/xm"
	"gitee.com/ywengineer/smart-kit/pkg/caches"
)

type Factory func(cp config.ChannelProperty) (inf.Verifier, error)

var factories = make(map[string]Factory)

var verifierCache caches.Cache[inf.Verifier]
var s sync.Once

func init() {
	s.Do(func() {
		var err error
		if verifierCache = caches.NewLocalCache[inf.Verifier](1000); verifierCache == nil {
			panic(err)
		}
		//
		RegisterFactory("rustore", vk.New)
		//
		RegisterFactory("huawei", hw.New)
		//
		RegisterFactory("xiaomi", xm.New)
	})
}

// RegisterFactory register a custom verifier factory
func RegisterFactory(name string, factory Factory) {
	factories[name] = factory
}

func FindVerifier(c config.Channel) (inf.Verifier, error) {
	// 从配置中获取通道属性
	cp, exists := config.Get().Channel[c.Code]
	if !exists {
		return nil, inf.ErrNoChannel
	}

	// 尝试从缓存中获取验证器
	var v inf.Verifier
	err := verifierCache.Get(cp.Validator, &v)
	if err == nil {
		return v, nil
	}

	// 获取对应的工厂函数
	factory, exists := factories[cp.Validator]
	if !exists {
		return nil, inf.ErrNoValidator
	}

	// 使用工厂函数创建验证器
	v, err = factory(cp)
	if err != nil {
		return nil, err
	}

	// 缓存验证器
	if err := verifierCache.Put(cp.Validator, v); err != nil {
		return nil, err
	}

	return v, nil
}