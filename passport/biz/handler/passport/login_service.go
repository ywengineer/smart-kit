// Code generated by hertz generator.

package passport

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/ywengineer/smart-kit/passport/biz/model/passport"
	"github.com/ywengineer/smart-kit/passport/internal"
	model2 "github.com/ywengineer/smart-kit/passport/internal/model"
	app2 "github.com/ywengineer/smart-kit/pkg/app"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"time"
)

// Login .
// @router /login [GET]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req passport.LoginReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, internal.ValidateErr(err))
		return
	}
	sCtx := ctx.Value(app2.ContextKeySmart).(app2.SmartContext)
	//
	c.JSON(consts.StatusOK, _login(ctx, sCtx, req.GetType(), req.GetId(), req.GetAccessToken(), ""))
}

func _login(ctx context.Context, sCtx app2.SmartContext, actType passport.AccountType, actId, token, refreshToken string) *app2.ApiResult {
	//
	bindKey := model2.GetBindCacheKey(actType.String(), actId)
	// query bind cache
	bkv, err := sCtx.Redis().JSONGet(ctx, bindKey).Result()
	var bind model2.PassportBinding
	//
	if errors.Is(err, redis.Nil) || len(bkv) == 0 { // query bind db
		r := sCtx.Rdb().
			WithContext(ctx).
			Where(&model2.PassportBinding{BindType: actType.String(), BindId: actId}).
			First(&bind)
		expire := 0
		if errors.Is(r.Error, gorm.ErrRecordNotFound) { // no data
			// cache null value for one minute
			bind.CreatedAt = time.Now()
			expire = 15 - bind.CreatedAt.Second()%10
		} else if r.Error != nil {
			hlog.Error("get data from rdb", zap.String("err", r.Error.Error()), zap.String("tag", "login_service"))
			return &internal.ErrRdb
		}
		if bs, err := sonic.Marshal(bind); err != nil { // json error
			return &internal.ErrJsonMarshal // stop
		} else if err = sCtx.Redis().JSONSet(ctx, bindKey, "$", bs).Err(); err != nil { // cache error
			hlog.Error("cache rdb object failed", zap.String("err", err.Error()), zap.String("tag", "login_service"))
			return &internal.ErrCache // stop
		} else if expire > 0 {
			sCtx.Redis().Expire(ctx, bindKey, time.Duration(expire)*time.Second)
		}
	} else if err != nil {
		hlog.Error("unreachable cache", zap.String("err", err.Error()), zap.String("tag", "login_service"))
		return &internal.ErrCache // stop
	} else if err = sonic.UnmarshalString(bkv, &bind); err != nil { // cache error
		hlog.Error("broken cache schema", zap.String("err", err.Error()), zap.String("tag", "login_service"))
		return &internal.ErrJsonUnmarshal // stop
	}
	//-------------------------------------- cache null --------------------------------------
	if bind.ID <= 0 {
		return &internal.ErrLoginTry // stop
	}
	//-------------------------------------- token match [Anonymous/EMail/Mobile] --------------------------------------
	if (actType == passport.AccountType_Anonymous || actType == passport.AccountType_Mobile || actType == passport.AccountType_EMail) && bind.AccessToken != token {
		return &internal.ErrInvalidToken
	} else { // update third platform token
		bind.AccessToken, bind.RefreshToken = token, refreshToken
		if ur := sCtx.Rdb().
			WithContext(ctx).
			Model(&bind).
			Select("AccessToken", "RefreshToken").Updates(model2.PassportBinding{AccessToken: token, RefreshToken: refreshToken}); ur.Error != nil || ur.RowsAffected == 0 {
			return &internal.ErrInvalidToken
		} else {
			_ = sCtx.Redis().JSONMerge(ctx, bindKey, "$", fmt.Sprintf(`{"access_token":"%s","refresh_token":"%s"}`, token, refreshToken))
		}
	}
	//-------------------------------------- after all process --------------------------------------
	if tk, _, err := sCtx.Jwt().TokenGenerator(map[string]interface{}{ // jwt token
		sCtx.Jwt().IdentityKey: bind.PassportId,
	}); err != nil {
		return &internal.ErrGenToken
	} else if bindTypes, err := sCtx.Redis().Get(ctx, internal.CacheKeyBoundTypes(bind.PassportId)).Result(); err != nil {
		return &internal.ErrCache
	} else {
		//
		pst := &model2.Passport{}
		if pstJson, err := sCtx.Redis().JSONGet(ctx, internal.CacheKeyPassport(bind.PassportId)).Result(); errors.Is(err, redis.Nil) || len(pstJson) <= 0 {
			// load from rdb
			pst.ID = bind.PassportId
			if fr := sCtx.Rdb().WithContext(ctx).Where(pst).First(pst); fr.Error != nil {
				hlog.Error("get passport data from rdb", zap.String("err", err.Error()), zap.String("tag", "login_service"))
				return &internal.ErrRdb
			}
			pstJson, _ = sonic.MarshalString(pst)
			_ = sCtx.Redis().JSONSet(ctx, internal.CacheKeyPassport(bind.PassportId), "$", pstJson)
		} else if err != nil {
			hlog.Error("get passport data from redis", zap.String("err", err.Error()), zap.String("tag", "login_service"))
			return &internal.ErrCache
		} else {
			_ = sonic.UnmarshalString(pstJson, pst)
		}
		//
		return app2.ApiOk(passport.LoginResp{
			PassportId: int64(bind.PassportId),
			Token:      tk,
			BrandNew:   false,
			Bounds: lo.Map(strings.Split(bindTypes, ","), func(item string, index int) passport.AccountType {
				ri, _ := passport.AccountTypeFromString(item)
				return ri
			}),
			CreateTime: pst.CreatedAt.Unix(),
		})
	}
}
