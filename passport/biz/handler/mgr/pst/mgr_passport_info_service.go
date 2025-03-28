// Code generated by hertz generator.

package pst

import (
	"context"
	"errors"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
	"github.com/ywengineer/smart-kit/passport/internal"
	"github.com/ywengineer/smart-kit/passport/internal/converter"
	model2 "github.com/ywengineer/smart-kit/passport/internal/model"
	app2 "github.com/ywengineer/smart-kit/pkg/apps"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	pst "github.com/ywengineer/smart-kit/passport/biz/model/mgr/pst"
)

// Detail .
// @router /mgr/passport/detail [GET]
func Detail(ctx context.Context, c *app.RequestContext) {
	var err error
	var req pst.MgrPassportDetailReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, internal.ValidateErr(err))
		return
	} //
	sCtx := ctx.Value(app2.ContextKeySmart).(app2.SmartContext)
	pstJson, err := sCtx.Redis().JSONGet(ctx, internal.CacheKeyPassport(uint(req.GetId()))).Result()
	if err == nil || len(pstJson) == 0 || errors.Is(err, redis.Nil) {
		var pstModel model2.Passport
		if len(pstJson) == 0 {
			if r := sCtx.Rdb().
				WithContext(ctx).
				Unscoped().
				Where("id = ?", req.GetId()).
				First(&pstModel); r.Error != nil && !errors.Is(r.Error, gorm.ErrRecordNotFound) {
				hlog.Error("failed to find passport", zap.String("err", r.Error.Error()), zap.Uint("data", pstModel.ID), zap.String("tag", "mgr_passport_info_service"))
				c.JSON(consts.StatusOK, internal.ErrRdb)
				return
			} else if errors.Is(r.Error, gorm.ErrRecordNotFound) {
				c.JSON(consts.StatusOK, internal.ErrUserNotFound)
				return
			}
		} else if err = sonic.UnmarshalString(pstJson, &pstModel); err != nil {
			hlog.Error("decode passport failed", zap.String("err", err.Error()), zap.String("data", pstJson), zap.String("tag", "mgr_passport_info_service"))
			c.JSON(consts.StatusOK, internal.ErrJsonUnmarshal)
			return
		}
		//
		var bounds []model2.PassportBinding
		if r := sCtx.Rdb().
			WithContext(ctx).
			Unscoped().
			Where(&model2.PassportBinding{PassportId: pstModel.ID}).
			Find(&bounds); r.Error != nil && !errors.Is(r.Error, gorm.ErrRecordNotFound) {
			hlog.Error("failed to find bounds data for passport", zap.String("err", r.Error.Error()), zap.Uint("data", pstModel.ID), zap.String("tag", "mgr_passport_info_service"))
			c.JSON(consts.StatusOK, internal.ErrRdb)
			return
		}
		//
		c.JSON(consts.StatusOK, app2.ApiOk(converter.ConvertPassport(&pstModel, &bounds)))
	} else {
		hlog.Error("unreachable cache", zap.String("err", err.Error()), zap.String("tag", "mgr_passport_info_service"))
		c.JSON(consts.StatusOK, internal.ErrCache)
	}
}
