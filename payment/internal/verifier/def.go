package verifier

import (
	"context"
	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
)

type Verifier interface {
	Verify(ctx context.Context, receipt string) (model.Purchase, error)
}

func FindVerifier(c config.Channel) (Verifier, error) {
	switch c.Code {
	//case config.ChannelAlipay:
	//	return &AlipayVerifier{}, nil
	}
	return nil, ErrNoChannel
}
