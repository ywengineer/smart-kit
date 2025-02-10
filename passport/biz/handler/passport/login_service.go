// Code generated by hertz generator.

package passport

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	passport "github.com/ywengineer/smart-kit/passport/biz/model/passport"
)

// Login .
// @router /login [GET]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req passport.LoginReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	//
	//sCtx := ctx.Value(pkg.ContextKeySmart).(*pkg.SmartContext)
	//
	resp := new(passport.LoginResp)

	c.JSON(consts.StatusOK, resp)
}
