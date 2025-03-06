// Code generated by hertz generator.

package passport

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/ywengineer/smart-kit/passport/pkg"
	"net/http"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _bindMw() []app.HandlerFunc {
	return []app.HandlerFunc{func(c context.Context, ctx *app.RequestContext) {
		v := c.Value(pkg.ContextKeySmart)
		if v == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		} else if sCtx, ok := v.(pkg.SmartContext); !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		} else {
			sCtx.TokenInterceptor()(c, ctx)
		}
	}}
}

func _loginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _registerMw() []app.HandlerFunc {
	return nil
}
