package model

import (
	"gitee.com/ywengineer/smart-kit/pkg/rdbs"
	"strings"
)

type PassportBinding struct {
	rdbs.Model
	PassportId   uint   `json:"passport_id" gorm:"index:idx_passport_id;comment:账号ID" redis:"passport_id"` //
	BindType     string `json:"bind_type" gorm:"index:idx_bind,unique;comment:绑定类型" redis:"bind_type"`     // 绑定类型：wx, qq, fb, apple, google
	BindId       string `json:"bind_id" gorm:"index:idx_bind,unique;comment:平台ID" redis:"bind_id"`         //
	AccessToken  string `json:"access_token"  gorm:"size:2000" redis:"access_token"`                       // 当绑定自己的邮件或手机号码时, 表示为密码
	RefreshToken string `json:"refresh_token" gorm:"size:2000" redis:"refresh_token"`
	SocialName   string `json:"social_name" redis:"social_name"`
	Gender       uint   `json:"gender" redis:"gender"`
	IconUrl      string `json:"icon_url" redis:"icon_url"`
}

func GetBindCacheKey(bindType, bindId string) string {
	return strings.Join([]string{"binding", bindType, bindId}, ":")
}
