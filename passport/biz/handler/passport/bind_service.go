// Code generated by hertz generator.

package passport

import (
	"context"
	"errors"
	"github.com/bsm/redislock"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/ywengineer/smart-kit/passport/internal"
	"github.com/ywengineer/smart-kit/passport/internal/model"
	"github.com/ywengineer/smart-kit/passport/pkg"
	"go.uber.org/zap"
	"strings"
	"time"

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
		c.JSON(consts.StatusBadRequest, internal.ValidateErr(err))
		return
	}
	// ano
	switch req.Type {
	case passport.AccountType_EMail, passport.AccountType_Mobile:
		c.JSON(consts.StatusNotImplemented, internal.ErrTodo)
		return
	default:
	}
	//
	sCtx := ctx.Value(pkg.ContextKeySmart).(pkg.SmartContext)
	passportId := uint(c.GetFloat64(sCtx.Jwt().IdentityKey))
	//----------------------------------------------- passport bind lock -----------------------------------------------
	lock, err := sCtx.LockMgr().Obtain(ctx, sCtx.GetPassportLockKey(passportId), time.Minute, &redislock.Options{
		Metadata:      "bind_service",
		RetryStrategy: redislock.NoRetry(),
	})
	if err != nil {
		hlog.Error("get lock err", zap.String("msg", err.Error()), zap.Uint("passportId", passportId), zap.String("tag", "bind_service"))
		c.JSON(consts.StatusLocked, internal.ErrDisLock)
		return
	}
	defer lock.Release(ctx)
	bindTypesCacheKey := internal.CacheKeyBoundTypes(passportId)
	// bound cache,
	bc, err := sCtx.Redis().Get(ctx, bindTypesCacheKey).Result()
	if errors.Is(err, redis.Nil) || len(bc) < 0 {
		// load
		var bounds []model.PassportBinding
		sCtx.Rdb().WithContext(ctx).Where(model.PassportBinding{PassportId: passportId}).Select("bind_type").Find(&bounds)
		bc = strings.Join(lo.Map(bounds, func(item model.PassportBinding, index int) string {
			return item.BindType
		}), ",")
		sCtx.Redis().Set(ctx, bindTypesCacheKey, bc, 0)
	}
	// already bound
	if strings.Contains(bc, req.GetType().String()) {
		c.JSON(consts.StatusOK, internal.ErrSameBound)
		return
	}
	// bind other
	bindKey := model.GetBindCacheKey(req.GetType().String(), req.GetBindId())
	// query bind cache
	if ok, err := sCtx.Redis().Exists(ctx, bindKey).Result(); err != nil && !errors.Is(err, redis.Nil) {
		c.JSON(consts.StatusOK, internal.ErrCache)
	} else if ok > 0 { // bind other
		c.JSON(consts.StatusOK, internal.ErrBoundOther)
	} else {
		// save rdb
		bindInfo := model.PassportBinding{
			PassportId:   passportId,
			BindType:     req.GetType().String(),
			BindId:       req.GetBindId(),
			AccessToken:  req.GetAccessToken(),
			RefreshToken: req.GetRefreshToken(),
			SocialName:   req.GetName(),
			Gender:       uint(req.GetGender()),
			IconUrl:      req.GetIconUrl(),
		}
		sCtx.Rdb().WithContext(ctx).Save(&bindInfo)
		//
		bc = strings.Join([]string{bc, req.Type.String()}, ",")
		if err = sCtx.Redis().Set(ctx, bindTypesCacheKey, bc, 0).Err(); err != nil {
			hlog.Error("cache last bound types failed", zap.String("err", err.Error()), zap.String("tag", "bind_service"))
			c.JSON(consts.StatusOK, internal.ErrCache)
		} else if bs, err := sonic.Marshal(bindInfo); err != nil { // json error
			c.JSON(consts.StatusOK, internal.ErrJsonMarshal)
		} else if err = sCtx.Redis().JSONSet(ctx, bindKey, "$", bs).Err(); err != nil { // cache error
			hlog.Error("cache new bind info failed", zap.String("err", err.Error()), zap.String("tag", "bind_service"))
			c.JSON(consts.StatusOK, internal.ErrCache)
		} else {
			c.JSON(consts.StatusOK, pkg.ApiOk(strings.Split(bc, ",")))
		}
	}
}
