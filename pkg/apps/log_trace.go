package apps

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/tracer/stats"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type tracerLog struct {
}

func (t *tracerLog) Start(ctx context.Context, c *app.RequestContext) context.Context {
	return ctx
}

func (t *tracerLog) Finish(ctx context.Context, c *app.RequestContext) {
	if ti := c.GetTraceInfo().Stats(); ti != nil {
		s, e := ti.GetEvent(stats.HTTPStart), ti.GetEvent(stats.HTTPFinish)
		if e != nil && s != nil {
			hlog.Infof("[Trace] %d [%s][%dms] [%s] [%s]", c.GetResponse().StatusCode(), consts.StatusMessage(c.GetResponse().StatusCode()), e.Time().Sub(s.Time()).Milliseconds(), c.Path(), e.Info())
		}
	}
}

func (t *tracerLog) statusText(status stats.Status) string {
	switch status {
	case stats.StatusError:
		return "ERROR"
	case stats.StatusWarn, stats.StatusInfo:
		return "OK"
	default:
		return "UNKNOWN"
	}
}
