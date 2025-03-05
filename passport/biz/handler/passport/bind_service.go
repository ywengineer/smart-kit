// Code generated by hertz generator.

package passport

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	passport "github.com/ywengineer/smart-kit/passport/biz/model/passport"
)

// Bind .
// @router /bind [GET]
func Bind(ctx context.Context, c *app.RequestContext) {
	var err error
	var req passport.BindReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, validateErr(err))
		return
	}

	resp := new(passport.LoginResp)

	c.JSON(consts.StatusOK, resp)
}
