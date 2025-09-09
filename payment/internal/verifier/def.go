package verifier

import (
	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/inf"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/vk"
)

type Factory func(cp config.ChannelProperty) (inf.Verifier, error)

var factories = make(map[string]Factory)

func init() {
	//
	RegisterFactory("rustore", func(cp config.ChannelProperty) (inf.Verifier, error) {
		return vk.NewRustore(vk.RustoreConfig{
			ClientID:     cp.ClientID,
			ClientSecret: cp.ClientSecret,
			IsSandbox:    cp.Sandbox,
			Apps:         cp.Apps,
		})
	})
	//
	RegisterFactory("huawei", func(cp config.ChannelProperty) (inf.Verifier, error) {
		return vk.NewRustore(vk.RustoreConfig{
			ClientID:     cp.ClientID,
			ClientSecret: cp.ClientSecret,
			IsSandbox:    cp.Sandbox,
			Apps:         cp.Apps,
		})
	})
	//
	RegisterFactory("xiaomi", func(cp config.ChannelProperty) (inf.Verifier, error) {
		return vk.NewRustore(vk.RustoreConfig{
			ClientID:     cp.ClientID,
			ClientSecret: cp.ClientSecret,
			IsSandbox:    cp.Sandbox,
			Apps:         cp.Apps,
		})
	})
}

// RegisterFactory register a custom verifier factory
func RegisterFactory(name string, factory Factory) {
	factories[name] = factory
}

func FindVerifier(c config.Channel) (inf.Verifier, error) {
	if cp, ok := config.Get().Channel[c.Code]; ok {
		if factory, ok := factories[cp.Validator]; ok {
			return factory(cp)
		} else {
			return nil, inf.ErrNoValidator
		}
	}
	return nil, inf.ErrNoChannel
}
