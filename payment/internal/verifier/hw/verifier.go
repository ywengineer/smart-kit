package hw

import (
	"context"

	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/payment/internal/verifier/inf"
	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
)

type huawei struct {
	config config.ChannelProperty
}

func (r *huawei) Verify(ctx context.Context, receipt string) (*model.Purchase, error) {
	//TODO implement me
	panic("implement me")
}

// New 初始化
func New(config config.ChannelProperty) (inf.Verifier, error) {
	hw := huawei{config: config}
	return &hw, nil
}
