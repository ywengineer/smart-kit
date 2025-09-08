package handler

import (
	"context"
	"errors"

	"gitee.com/ywengineer/smart-kit/payment/internal/queue"
	"gitee.com/ywengineer/smart-kit/payment/internal/services"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hibiken/asynq"
)

func Test(ctx context.Context, task *asynq.Task) error {
	var q queue.TestAsynqQueue
	if err := sonic.Unmarshal(task.Payload(), &q); err != nil {
		return err
	}
	hlog.CtxInfof(ctx, "test queue, payload = %s", string(task.Payload()))
	if q.Flag {
		return errors.New("test failed task")
	}
	return nil
}

func HandlePurchaseNotify(ctx context.Context, t *asynq.Task) error {
	var payload queue.PurchaseNotifyPayload
	if err := sonic.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return services.Notify(ctx, payload)
}
