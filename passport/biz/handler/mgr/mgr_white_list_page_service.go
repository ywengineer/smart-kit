// Code generated by hertz generator.

package mgr

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/ywengineer/smart-kit/passport/internal"
	"github.com/ywengineer/smart-kit/passport/internal/converter"
	model2 "github.com/ywengineer/smart-kit/passport/internal/model"
	app2 "github.com/ywengineer/smart-kit/pkg/app"
	"github.com/ywengineer/smart-kit/pkg/sql"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	mgr "github.com/ywengineer/smart-kit/passport/biz/model/mgr"
)

// Page .
// @router /mgr/white-list/page [GET]
func Page(ctx context.Context, c *app.RequestContext) {
	var err error
	var req mgr.WhiteListPageReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusBadRequest, internal.ValidateErr(err))
		return
	}
	//
	sCtx := ctx.Value(app2.ContextKeySmart).(app2.SmartContext)
	//
	if req.PassportId > 0 {
		var user model2.WhiteList
		//
		if stmt := sCtx.Rdb().WithContext(ctx).
			Unscoped().
			Where(&model2.WhiteList{Passport: uint(req.GetPassportId())}).
			First(&user); stmt.Error != nil && !errors.Is(stmt.Error, gorm.ErrRecordNotFound) {
			hlog.Error("paginator error", zap.Any("data", req), zap.String("err", stmt.Error.Error()), zap.String("tag", "white_list_page_service"))
			c.JSON(consts.StatusOK, internal.ErrRdb)
		} else {
			c.JSON(consts.StatusOK, app2.ApiOk(&mgr.WhiteListPageRes{
				Page:     1,
				PageSize: req.GetPageSize(),
				Total:    1,
				MaxPage:  1,
				Data:     []*mgr.WhiteListData{converter.WhiteList(user)},
			}))
		}
	} else {
		var users []model2.WhiteList
		// extend query before paginating
		stmt := sCtx.Rdb().WithContext(ctx).Model(&users).Order("id desc")
		// with pagination
		if p, err := sql.Paginate[model2.WhiteList](ctx, stmt, req.GetPageNo(), req.GetPageSize(), &users, converter.WhiteList); err != nil {
			hlog.Error("paginator error", zap.Any("data", req), zap.String("err", err.Error()), zap.String("tag", "white_list_page_service"))
			c.JSON(consts.StatusOK, internal.ErrRdb)
		} else {
			c.JSON(consts.StatusOK, app2.ApiOk(&mgr.WhiteListPageRes{
				Page:     req.GetPageNo(),
				PageSize: req.GetPageSize(),
				Total:    p.Total,
				MaxPage:  p.MaxPage,
				Data:     p.Data,
			}))
		}
	}
}
