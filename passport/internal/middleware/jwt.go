package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/jwt"
	"github.com/ywengineer/smart-kit/pkg/apps"
)

func Jwt() []app.HandlerFunc {
	return []app.HandlerFunc{func(c context.Context, ctx *app.RequestContext) {
		v := apps.GetContext(c)
		if v == nil {
			ctx.AbortWithStatus(consts.StatusServiceUnavailable)
		} else {
			v.TokenInterceptor()(c, ctx)
		}
	}}
}

func JwtWithValidate(f func(data jwt.MapClaims) bool) []app.HandlerFunc {
	return []app.HandlerFunc{
		Jwt()[0],
		func(c context.Context, ctx *app.RequestContext) {
			if f(jwt.ExtractClaims(c, ctx)) {
				ctx.Next(c)
			} else {
				ctx.AbortWithStatus(consts.StatusUnauthorized)
			}
		},
	}
}
