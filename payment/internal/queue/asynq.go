package queue

import (
	"context"
	"sync"

	"gitee.com/ywengineer/smart-kit/payment/internal/config"
	"gitee.com/ywengineer/smart-kit/pkg/apps"
	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hibiken/asynq"
)

var o sync.Once
var so sync.Once

var cli *asynq.Client
var srv *asynq.Server

func InitQueue(ctx context.Context, sCtx apps.SmartContext, cfg config.Queue, handlers map[TaskType]asynq.HandlerFunc) {
	o.Do(func() {
		ctx := context.WithValue(ctx, apps.ContextKeySmart, sCtx)
		cli = asynq.NewClientFromRedisClient(sCtx.Redis())
		srv = asynq.NewServerFromRedisClient(sCtx.Redis(),
			asynq.Config{
				Concurrency: cfg.Workers,
				Queues:      cfg.Queues,
				BaseContext: func() context.Context {
					return ctx
				},
				Logger: logk.NewSLogger("./logs/queue.log", 10, 10, 7, hlog.LevelInfo),
			})
		// mux maps a type to a handler
		mux := asynq.NewServeMux()
		for t, h := range handlers {
			mux.Handle(string(t), h)
		}
		hlog.Infof("start queue service: %v", srv.Start(mux))
	})
}

func Shutdown() {
	so.Do(func() {
		hlog.Infof("close queue client: %v", cli.Close())
		srv.Stop()
		srv.Shutdown()
	})
}
