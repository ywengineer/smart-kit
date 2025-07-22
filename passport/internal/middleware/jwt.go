package middleware

import (
	"context"
	"gitee.com/ywengineer/smart-kit/passport/internal"
	"gitee.com/ywengineer/smart-kit/pkg/apps"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/jwt"
)

type TokenValidator func(data jwt.MapClaims) bool

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

func JwtWithValidate(f TokenValidator) []app.HandlerFunc {
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

func IsUserMatch(tp internal.UserType) TokenValidator {
	return func(data jwt.MapClaims) bool {
		if ut, ok := data[internal.TokenKeyUserType]; !ok {
			return false
		} else if internal.UserType(ut.(float64)) != tp {
			return false
		} else {
			return true
		}
	}
}
