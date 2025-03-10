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
	"github.com/ywengineer/smart-kit/passport/pkg"
	"github.com/ywengineer/smart-kit/passport/pkg/model"
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
		c.JSON(consts.StatusBadRequest, validateErr(err))
		return
	}
	sCtx := ctx.Value(pkg.ContextKeySmart).(pkg.SmartContext)
	//
	c.JSON(consts.StatusOK, _login(ctx, sCtx, req.GetType(), req.GetId(), req.GetAccessToken(), ""))
}

func _login(ctx context.Context, sCtx pkg.SmartContext, actType passport.AccountType, actId, token, refreshToken string) *pkg.ApiResult {
	//
	bindKey := model.GetBindCacheKey(actType.String(), actId)
	// query bind cache
	bkv, err := sCtx.Redis().JSONGet(ctx, bindKey).Result()
	var bind model.PassportBinding
	//
	if errors.Is(err, redis.Nil) || len(bkv) == 0 { // query bind db
		r := sCtx.Rdb().
			WithContext(ctx).
			Where(&model.PassportBinding{BindType: actType.String(), BindId: actId}).
			First(&bind)
		expire := 0
		if errors.Is(r.Error, gorm.ErrRecordNotFound) { // no data
			// cache null value for one minute
			bind.CreatedAt = time.Now()
			expire = 15 - bind.CreatedAt.Second()%10
		} else if r.Error != nil {
			hlog.Error("get data from rdb", zap.String("err", r.Error.Error()), zap.String("tag", "login_service"))
			return &ErrRdb
		}
		if bs, err := sonic.Marshal(bind); err != nil { // json error
			return &ErrJsonMarshal // stop
		} else if err = sCtx.Redis().JSONSet(ctx, bindKey, "$", bs).Err(); err != nil { // cache error
			hlog.Error("cache rdb object failed", zap.String("err", err.Error()), zap.String("tag", "login_service"))
			return &ErrCache // stop
		} else if expire > 0 {
			sCtx.Redis().Expire(ctx, bindKey, time.Duration(expire)*time.Second)
		}
	} else if err != nil {
		hlog.Error("unreachable cache", zap.String("err", err.Error()), zap.String("tag", "login_service"))
		return &ErrCache // stop
	} else if err = sonic.UnmarshalString(bkv, &bind); err != nil { // cache error
		hlog.Error("broken cache schema", zap.String("err", err.Error()), zap.String("tag", "login_service"))
		return &ErrJsonUnmarshal // stop
	}
	//-------------------------------------- cache null --------------------------------------
	if bind.ID <= 0 {
		return &ErrLoginTry // stop
	}
	//-------------------------------------- token match [Anonymous/EMail/Mobile] --------------------------------------
	if (actType == passport.AccountType_Anonymous || actType == passport.AccountType_Mobile || actType == passport.AccountType_EMail) && bind.AccessToken != token {
		return &ErrInvalidToken
	} else { // update third platform token
		bind.AccessToken, bind.RefreshToken = token, refreshToken
		if ur := sCtx.Rdb().
			WithContext(ctx).
			Model(&bind).
			Select("AccessToken", "RefreshToken").Updates(model.PassportBinding{AccessToken: token, RefreshToken: refreshToken}); ur.Error != nil || ur.RowsAffected == 0 {
			return &ErrInvalidToken
		} else {
			_ = sCtx.Redis().JSONMerge(ctx, bindKey, "$", fmt.Sprintf(`{"access_token":"%s","refresh_token":"%s"}`, token, refreshToken))
		}
	}
	//-------------------------------------- after all process --------------------------------------
	if tk, _, err := sCtx.Jwt().TokenGenerator(map[string]interface{}{ // jwt token
		sCtx.Jwt().IdentityKey: bind.PassportId,
	}); err != nil {
		return &ErrGenToken
	} else if bindTypes, err := sCtx.Redis().Get(ctx, cacheKeyBoundTypes(bind.PassportId)).Result(); err != nil {
		return &ErrCache
	} else {
		//
		pst := &model.Passport{}
		if pstJson, err := sCtx.Redis().JSONGet(ctx, cacheKeyPassport(bind.PassportId), "$").Result(); errors.Is(err, redis.Nil) || len(pstJson) <= 0 {
			// load from rdb
			pst.ID = bind.PassportId
			sCtx.Rdb().WithContext(ctx).Where(pst).First(pst)
		} else if err != nil {
			hlog.Error("get passport data from redis", zap.String("err", err.Error()), zap.String("tag", "login_service"))
			return &ErrCache
		} else {
			_ = sonic.UnmarshalString(pstJson, pst)
		}
		//
		ok := pkg.ApiOk(passport.LoginResp{
			PassportId: int64(bind.PassportId),
			Token:      tk,
			BrandNew:   false,
			Bounds: lo.Map(strings.Split(bindTypes, ","), func(item string, index int) passport.AccountType {
				ri, _ := passport.AccountTypeFromString(item)
				return ri
			}),
			CreateTime: pst.CreatedAt.Unix(),
		})
		return &ok
	}
}
