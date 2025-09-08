package handler

import (
	"context"
	"gitee.com/ywengineer/smart-kit/payment/internal/queue"
	"gitee.com/ywengineer/smart-kit/payment/internal/service"
	"github.com/bytedance/sonic"
	"github.com/hibiken/asynq"
)

func HandlePurchaseNotify(ctx context.Context, t *asynq.Task) error {
	var payload queue.PurchaseNotifyPayload
	if err := sonic.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return service.Notify(ctx, payload)
}
