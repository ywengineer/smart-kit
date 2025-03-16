package converter

import (
	"github.com/ywengineer/smart-kit/passport/biz/model/mgr"
	model2 "github.com/ywengineer/smart-kit/passport/internal/model"
)

func WhiteList(i model2.WhiteList) *mgr.WhiteListData {
	return &mgr.WhiteListData{
		Id:         int64(i.ID),
		CreateAt:   i.CreatedAt.Unix(),
		UpdateAt:   i.UpdatedAt.Unix(),
		DeleteAt:   i.DeletedAt.Unix(),
		PassportId: int64(i.Passport),
	}
}
