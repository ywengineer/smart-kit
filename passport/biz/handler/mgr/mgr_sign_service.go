// Code generated by hertz generator.

package mgr

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
	"github.com/ywengineer/smart-kit/passport/internal"
	model2 "github.com/ywengineer/smart-kit/passport/internal/model"
	"github.com/ywengineer/smart-kit/passport/pkg"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	mgr "github.com/ywengineer/smart-kit/passport/biz/model/mgr"
)

// Sign .
// @router /mgr/sign [POST]
func Sign(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mgr.MgrSignReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, internal.ValidateErr(err))
		return
	}
	c.JSON(consts.StatusOK, mgrSign(ctx, c, &req))
}

func mgrSign(ctx context.Context, c *app.RequestContext, req *mgr.MgrSignReq) *pkg.ApiResult {
	//
	sCtx := ctx.Value(pkg.ContextKeySmart).(pkg.SmartContext)
	//
	mgrKey := model2.GetMgrCacheKey(req.GetAccount())
	// query bind cache
	bkv, err := sCtx.Redis().JSONGet(ctx, mgrKey).Result()
	var bind model2.MgrUser
	//
	if errors.Is(err, redis.Nil) || len(bkv) == 0 { // query bind db
		r := sCtx.Rdb().
			WithContext(ctx).
			Where(&model2.MgrUser{Account: req.GetAccount()}).
			First(&bind)
		expire := 0
		if errors.Is(r.Error, gorm.ErrRecordNotFound) { // no data
			// cache null value for one minute
			bind.CreatedAt = time.Now()
			expire = 15 - bind.CreatedAt.Second()%10
		} else if r.Error != nil {
			hlog.Error("get data from rdb", zap.String("err", r.Error.Error()), zap.String("tag", "mgr_sign_service"))
			return &internal.ErrRdb
		}
		if bs, err := sonic.Marshal(bind); err != nil { // json error
			return &internal.ErrJsonMarshal // stop
		} else if err = sCtx.Redis().JSONSet(ctx, mgrKey, "$", bs).Err(); err != nil { // cache error
			hlog.Error("cache rdb object failed", zap.String("err", err.Error()), zap.String("tag", "mgr_sign_service"))
			return &internal.ErrCache // stop
		} else if expire > 0 {
			sCtx.Redis().Expire(ctx, mgrKey, time.Duration(expire)*time.Second)
			return &internal.ErrUserNotFound
		}
	} else if err != nil {
		hlog.Error("unreachable cache", zap.String("err", err.Error()), zap.String("tag", "mgr_sign_service"))
		return &internal.ErrCache // stop
	} else if err = sonic.UnmarshalString(bkv, &bind); err != nil { // cache error
		hlog.Error("broken cache schema", zap.String("err", err.Error()), zap.String("tag", "mgr_sign_service"))
		return &internal.ErrJsonUnmarshal // stop
	}
	//-------------------------------------- after all process --------------------------------------
	if bind.Password != req.GetPassword() {
		return &internal.ErrPassword
	} else if tk, _, err := sCtx.Jwt().TokenGenerator(map[string]interface{}{ // jwt token
		sCtx.Jwt().IdentityKey: bind.ID,
	}); err != nil {
		return &internal.ErrGenToken
	} else {
		//
		return pkg.ApiOk(&mgr.MgrSignRes{
			Id:    int64(bind.ID),
			Act:   req.GetAccount(),
			Name:  bind.Name,
			Dept:  int64(bind.DeptNo),
			Title: bind.Title,
			Token: tk,
		})
	}
}
