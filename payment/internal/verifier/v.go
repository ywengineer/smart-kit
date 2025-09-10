package verifier

import (
	"sync"

	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/hw"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/inf"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/vk"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/xm"
	"gitee.com/ywengineer/smart-kit/pkg/caches"
	"github.com/dgraph-io/ristretto/v2"
)

type Factory func(cp config.ChannelProperty) (inf.Verifier, error)

var factories = make(map[string]Factory)

var verifierCache *ristretto.Cache[string, inf.Verifier]
var s sync.Once

func init() {
	s.Do(func() {
		var err error
		if verifierCache, err = caches.NewCache[inf.Verifier](1000); err != nil {
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
	if cp, ok := config.Get().Channel[c.Code]; ok {
		if v, ok := verifierCache.Get(cp.Validator); ok {
			return v, nil
		} else if factory, ok := factories[cp.Validator]; ok {
			if v, err := factory(cp); err == nil {
				verifierCache.Set(cp.Validator, v, 1)
				return v, nil
			} else {
				return nil, err
			}
		} else {
			return nil, inf.ErrNoValidator
		}
	}
	return nil, inf.ErrNoChannel
}
