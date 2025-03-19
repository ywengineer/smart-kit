package converter

import (
	"github.com/samber/lo"
	"github.com/ywengineer/smart-kit/passport/biz/model/mgr/pst"
	model2 "github.com/ywengineer/smart-kit/passport/internal/model"
)

func ConvertPassport(pstModel *model2.Passport, bounds *[]model2.PassportBinding) *pst.PassportData {
	pd := &pst.PassportData{
		Id:         int64(pstModel.ID),
		CreateAt:   pstModel.CreatedAt.Unix(),
		UpdateAt:   pstModel.UpdatedAt.Unix(),
		DeleteAt:   pstModel.DeletedAt.Unix(),
		DeviceId:   pstModel.DeviceId,
		Adid:       pstModel.Adid,
		SystemType: pstModel.SystemType,
		Locale:     pstModel.Locale,
		Extra:      pstModel.Extra.String(),
	}
	if bounds != nil && len(*bounds) > 0 {
		pd.Bounds = lo.Map(*bounds, func(item model2.PassportBinding, index int) *pst.PassportBindData {
			return &pst.PassportBindData{
				Id:         int64(item.ID),
				CreateAt:   item.CreatedAt.Unix(),
				UpdateAt:   item.UpdatedAt.Unix(),
				DeleteAt:   item.DeletedAt.Unix(),
				BindType:   item.BindType,
				BindId:     item.BindId,
				Token:      item.AccessToken,
				SocialName: item.SocialName,
				Gender:     int64(item.Gender),
				IconUrl:    item.IconUrl,
			}
		})
	}
	return pd
}
