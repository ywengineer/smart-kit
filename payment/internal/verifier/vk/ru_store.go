package vk

import (
	"context"
	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
)

type ruStore struct {
}

func (r ruStore) Verify(ctx context.Context, receipt string) (model.Purchase, error) {
	panic("implement me")
}
