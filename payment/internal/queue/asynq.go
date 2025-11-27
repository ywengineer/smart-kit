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

func Setup(ctx context.Context, sCtx apps.SmartContext, cfg config.Queue, handlers map[TaskType]asynq.HandlerFunc) {
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
				Logger: logk.NewSLogger("./logs/queue.log", logk.WithLevel(logk.InfoLevel)),
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
		srv.Stop()
		srv.Shutdown()
		_ = cli.Close()
	})
}
