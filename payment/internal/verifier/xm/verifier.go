package xm

import (
	"context"

	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/inf"
	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
)

type xiaomi struct {
	config config.ChannelProperty
}

// New 初始化
func New(config config.ChannelProperty) (inf.Verifier, error) {
	hw := xiaomi{config: config}
	return &hw, nil
}

func (r *xiaomi) Verify(ctx context.Context, receipt string) (*model.Purchase, error) {
	panic("implement me")
}
