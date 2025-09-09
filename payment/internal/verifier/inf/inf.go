package inf

import (
	"context"

	"gitee.com/ywengineer/smart-kit/payment/pkg/model"
)

type Verifier interface {
	Verify(ctx context.Context, receipt string) (*model.Purchase, error)
}
