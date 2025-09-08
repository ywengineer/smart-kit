package handler

import (
	"context"
	"gitee.com/ywengineer/smart-kit/payment/internal/queue"
	"gitee.com/ywengineer/smart-kit/payment/internal/service"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hibiken/asynq"
)

func Test(ctx context.Context, task *asynq.Task) error {
	hlog.CtxInfof(ctx, "test queue, payload = %s", string(task.Payload()))
	return nil
}

func HandlePurchaseNotify(ctx context.Context, t *asynq.Task) error {
	var payload queue.PurchaseNotifyPayload
	if err := sonic.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return service.Notify(ctx, payload)
}
