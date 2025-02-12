// Code generated by hertz generator.

package passport

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/redis/go-redis/v9"
	"github.com/ywengineer/smart-kit/passport/biz/model/passport"
	"github.com/ywengineer/smart-kit/passport/pkg"
	"github.com/ywengineer/smart-kit/passport/pkg/model"
	"gorm.io/gorm"
	"math/rand/v2"
	"time"
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
	sCtx := ctx.Value(pkg.ContextKeySmart).(pkg.SmartContext)
	bindKey := model.GetBindCacheKey(req.GetType().String(), req.GetId())
	// query bind cache
	cv, err := sCtx.Redis().Get(ctx, bindKey).Result()
	var bind model.PassportBinding
	//
	if errors.Is(err, redis.Nil) { // query bind db
		r := sCtx.Rdb().
			WithContext(ctx).
			Where(&model.PassportBinding{BindType: req.GetType().String(), BindId: req.GetId()}).
			First(&bind)
		ex := 5 * 24 * time.Hour
		if errors.Is(r.Error, gorm.ErrRecordNotFound) { // no data
			// cache null value for one minute
			cv, ex = "", time.Minute
		} else { // cache result
			cv, _ = sonic.MarshalString(bind)
			ex += time.Duration(rand.Int64N(60) * int64(time.Minute))
		}
		// cache error
		if err = sCtx.Redis().SetNX(ctx, bindKey, cv, ex).Err(); err != nil {
			c.JSON(consts.StatusOK, pkg.ApiError(err.Error()))
			return // stop
		}
	} else if cv == "" { // cache null
		c.JSON(consts.StatusOK, pkg.ApiError("not found"))
		return // stop
	} else if err = sonic.UnmarshalString(cv, &bind); err != nil { // cache error
		c.JSON(consts.StatusOK, pkg.ApiError(err.Error()))
		return // stop
	} else {
		// obtain binding data success: ignore
	}
	//-------------------------------------- token match --------------------------------------
	if bind.AccessToken != req.GetAccessToken() {
		c.JSON(consts.StatusOK, pkg.ApiError("invalid.token"))
	} else if tk, _, err := sCtx.Jwt().TokenGenerator(map[string]interface{}{ // jwt token
		"id": bind.PassportId,
	}); err != nil {
		c.JSON(consts.StatusOK, pkg.ApiError(err.Error()))
	} else {
		c.JSON(consts.StatusOK, pkg.ApiOk(passport.LoginResp{
			PassportId: int64(bind.PassportId),
			Token:      tk,
			BrandNew:   false,
			CreateTime: time.Now().Unix(),
		}))
	}
}
