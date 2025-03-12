package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/ywengineer/smart-kit/passport/pkg"
)

func Jwt() []app.HandlerFunc {
	return []app.HandlerFunc{func(c context.Context, ctx *app.RequestContext) {
		v := c.Value(pkg.ContextKeySmart)
		if v == nil {
			ctx.AbortWithStatus(consts.StatusOK)
		} else if sCtx, ok := v.(pkg.SmartContext); !ok {
			ctx.AbortWithStatus(consts.StatusUnauthorized)
		} else {
			sCtx.TokenInterceptor()(c, ctx)
		}
	}}
}
