package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	app2 "github.com/ywengineer/smart-kit/pkg/app"
)

func Jwt() []app.HandlerFunc {
	return []app.HandlerFunc{func(c context.Context, ctx *app.RequestContext) {
		v := c.Value(app2.ContextKeySmart)
		if v == nil {
			ctx.AbortWithStatus(consts.StatusOK)
		} else if sCtx, ok := v.(app2.SmartContext); !ok {
			ctx.AbortWithStatus(consts.StatusUnauthorized)
		} else {
			sCtx.TokenInterceptor()(c, ctx)
		}
	}}
}
