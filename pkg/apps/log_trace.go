package apps

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/tracer/stats"
)

type tracerLog struct {
}

func (t *tracerLog) Start(ctx context.Context, c *app.RequestContext) context.Context {
	return ctx
}

func (t *tracerLog) Finish(ctx context.Context, c *app.RequestContext) {
	ti := c.GetTraceInfo().Stats()
	s, e := ti.GetEvent(stats.HTTPStart), ti.GetEvent(stats.HTTPFinish)
	hlog.Infof("[Trace] [%dms] [%s]", e.Time().Sub(s.Time()).Milliseconds(), c.Path())
}
