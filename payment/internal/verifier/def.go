package verifier

import (
	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/inf"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/vk"
)

type Factory func(cp config.ChannelProperty) (inf.Verifier, error)

var factories = map[string]Factory{
	"rustore": func(cp config.ChannelProperty) (inf.Verifier, error) {
		return vk.NewRustore(vk.RustoreConfig{
			ClientID:     cp.ClientID,
			ClientSecret: cp.ClientSecret,
			IsSandbox:    cp.Sandbox,
			Apps:         cp.Apps,
		})
	},
	"huawei": func(cp config.ChannelProperty) (inf.Verifier, error) {
		return nil, nil
	},
	"xiaomi": func(cp config.ChannelProperty) (inf.Verifier, error) {
		return nil, nil
	},
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
