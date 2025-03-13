// Code generated by hertz generator.

package passport

import (
	"context"
	"errors"
	"github.com/bsm/redislock"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/google/uuid"
	passport "github.com/ywengineer/smart-kit/passport/biz/model/passport"
	"github.com/ywengineer/smart-kit/passport/internal"
	model2 "github.com/ywengineer/smart-kit/passport/internal/model"
	"github.com/ywengineer/smart-kit/passport/pkg"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"time"
)

// Register .
// @router /register [GET]
func Register(ctx context.Context, c *app.RequestContext) {
	var req passport.RegisterReq
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, internal.ValidateErr(err))
		return
	}
	//
	sCtx := ctx.Value(pkg.ContextKeySmart).(pkg.SmartContext)
	bindKey := model2.GetBindCacheKey(req.GetType().String(), req.GetId())
	// ano
	switch req.Type {
	case passport.AccountType_EMail, passport.AccountType_Mobile:
		c.JSON(consts.StatusNotImplemented, internal.ErrTodo)
		return
	case passport.AccountType_Anonymous:
		// continue
	default: // other platform
		// exists
		if exists, err := sCtx.Redis().Exists(ctx, bindKey).Result(); err != nil {
			hlog.Error("exists check", zap.String("msg", err.Error()), zap.String("deviceId", req.DeviceId), zap.String("tag", "register_service"))
			c.JSON(consts.StatusInternalServerError, internal.ErrCache)
			return
		} else if exists > 0 { // already bind to passport, go to log in service
			c.JSON(consts.StatusOK, _login(ctx, sCtx, req.GetType(), req.GetId(), req.GetAccessToken(), req.GetRefreshToken()))
			return
		}
	}
	//----------------------------------------------- device lock -----------------------------------------------
	lock, err := sCtx.LockMgr().Obtain(ctx, sCtx.GetDeviceLockKey(req.GetDeviceId()), time.Minute, &redislock.Options{
		Metadata:      "register_service",
		RetryStrategy: redislock.NoRetry(),
	})
	if err != nil {
		hlog.Error("get lock err", zap.String("msg", err.Error()), zap.String("deviceId", req.DeviceId), zap.String("tag", "register_service"))
		c.JSON(consts.StatusLocked, internal.ErrDisLock)
		return
	}
	defer lock.Release(ctx)
	cntKey := "register:" + req.DeviceId
	//----------------------------------------------- max per device -----------------------------------------------
	cntNow, err := sCtx.Redis().IncrBy(ctx, cntKey, 0).Result()
	if err != nil {
		hlog.Error("get incr 0 err", zap.String("msg", err.Error()), zap.String("deviceId", req.DeviceId), zap.String("tag", "register_service"))
		c.JSON(consts.StatusInternalServerError, internal.ErrCache)
		return
	} else if cntNow >= 3 {
		hlog.Info("reach max per device", zap.String("deviceId", req.DeviceId), zap.String("tag", "register_service"))
		c.JSON(consts.StatusOK, internal.ErrMaxPerDevice)
		return
	}
	//----------------------------------------------- exists type and id? -----------------------------------------------
	if req.Type == passport.AccountType_Anonymous { // gen random id and reset bind cache key
		req.Id = strings.ToLower(strings.ReplaceAll(uuid.New().String(), "-", ""))
		req.AccessToken = uuid.New().String()
		bindKey = model2.GetBindCacheKey(req.GetType().String(), req.GetId())
	}
	var bind model2.PassportBinding
	if exists, err := sCtx.Redis().Exists(ctx, bindKey).Result(); err != nil {
		hlog.Error("exists check", zap.String("msg", err.Error()), zap.String("deviceId", req.DeviceId), zap.String("tag", "register_service"))
		c.JSON(consts.StatusInternalServerError, internal.ErrCache)
		return
	} else if exists > 0 {
		c.JSON(consts.StatusOK, internal.ErrBoundOther) // already bind to other passport
		return
	} else { // load from db
		r := sCtx.Rdb().
			WithContext(ctx).
			Where(&model2.PassportBinding{BindType: req.GetType().String(), BindId: req.GetId()}).
			First(&bind)
		// rdb error
		if r.Error != nil && !errors.Is(r.Error, gorm.ErrRecordNotFound) {
			hlog.Error("get data from rdb", zap.String("err", r.Error.Error()), zap.String("tag", "register_service"))
			c.JSON(consts.StatusOK, internal.ErrRdb)
			return
		}
		if bind.PassportId > 0 { // already bound
			//  cache
			if bs, err := sonic.Marshal(bind); err != nil { // json error
				c.JSON(consts.StatusOK, internal.ErrJsonMarshal)
				return // stop
			} else if err = sCtx.Redis().JSONSet(ctx, bindKey, "$", bs).Err(); err != nil { // cache error
				hlog.Error("cache rdb object failed", zap.String("err", err.Error()), zap.String("tag", "register_service"))
				c.JSON(consts.StatusOK, internal.ErrCache)
				return // stop
			} else {
				c.JSON(consts.StatusOK, internal.ErrBoundOther)
				return // stop
			}
		}
	}
	//----------------------------------------------- insert passport and binding -----------------------------------------------
	deviceBytes, _ := sonic.Marshal(req.GetDeviceInfo())
	pst := model2.Passport{
		DeviceId:   req.GetDeviceId(),
		Adid:       req.GetAdid(),
		SystemType: req.GetDeviceInfo()[pkg.Os],
		Locale:     req.GetDeviceInfo()[pkg.Locale],
		Extra:      deviceBytes,
	}
	if err = sCtx.Rdb().Transaction(func(pstBind *model2.PassportBinding) func(tx *gorm.DB) error {
		return func(tx *gorm.DB) error {
			if err := tx.Create(&pst).Error; err != nil { // return any error will roll back
				return err
			}
			*pstBind = model2.PassportBinding{
				PassportId:   pst.ID,
				BindType:     req.GetType().String(),
				BindId:       req.GetId(),
				AccessToken:  req.GetAccessToken(),
				RefreshToken: req.GetRefreshToken(),
				SocialName:   req.GetName(),
				Gender:       uint(req.GetGender()),
				IconUrl:      req.GetIconUrl(),
			}
			if err := tx.Create(pstBind).Error; err != nil {
				return err
			}
			return nil
		}
	}(&bind)); err != nil {
		hlog.Error("save rdb error", zap.String("msg", err.Error()), zap.String("deviceId", req.DeviceId), zap.String("tag", "register_service"))
		c.JSON(consts.StatusOK, internal.ErrRegisterFail)
		return
	}
	//
	sCtx.Redis().Incr(ctx, cntKey)
	//----------------------------------------------- finish -----------------------------------------------
	// cache
	if _, err := sCtx.Redis().Set(ctx, internal.CacheKeyBoundTypes(pst.ID), req.GetType().String(), 0).Result(); err != nil {
		hlog.Error("cache bind types failed", zap.String("err", err.Error()), zap.String("tag", "register_service"))
		c.JSON(consts.StatusOK, internal.ErrCache)
	} else if bs, err := sonic.Marshal(bind); err != nil { // json error
		c.JSON(consts.StatusOK, internal.ErrJsonMarshal)
	} else if err = sCtx.Redis().JSONSet(ctx, bindKey, "$", bs).Err(); err != nil { // cache error
		hlog.Error("cache new rdb object failed", zap.String("err", err.Error()), zap.String("tag", "register_service"))
		c.JSON(consts.StatusOK, internal.ErrCache)
	} else if pstJsonStr, err := sonic.Marshal(pst); err != nil { // json error
		c.JSON(consts.StatusOK, internal.ErrJsonMarshal)
	} else if err = sCtx.Redis().JSONSet(ctx, internal.CacheKeyPassport(pst.ID), "$", pstJsonStr).Err(); err != nil { // cache error
		hlog.Error("cache passport json object failed", zap.String("err", err.Error()), zap.String("tag", "register_service"))
		c.JSON(consts.StatusOK, internal.ErrCache)
	} else if tk, _, err := sCtx.Jwt().TokenGenerator(map[string]interface{}{ // jwt token
		sCtx.Jwt().IdentityKey: pst.ID,
	}); err != nil {
		hlog.Error("gen token error", zap.String("msg", err.Error()), zap.String("deviceId", req.DeviceId), zap.String("tag", "register_service"))
		c.JSON(consts.StatusOK, internal.ErrGenToken)
	} else {
		c.JSON(consts.StatusOK, pkg.ApiOk(passport.LoginResp{
			PassportId: int64(pst.ID),
			Token:      tk,
			BrandNew:   true,
			Bounds:     []passport.AccountType{req.Type},
			CreateTime: pst.CreatedAt.Unix(),
		}))
	}
}
