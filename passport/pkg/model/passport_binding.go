package model

import (
	"gorm.io/gorm"
	"strings"
)

type PassportBinding struct {
	gorm.Model
	PassportId   uint   `json:"passport_id" gorm:"index:idx_passport_id;comment:账号ID"` //
	BindType     string `json:"bind_type" gorm:"index:idx_bind,unique;comment:绑定类型"`   // 绑定类型：wx, qq, fb, apple, google
	BindId       string `json:"bind_id" gorm:"index:idx_bind,unique;comment:平台ID"`     //
	AccessToken  string `json:"access_token"  gorm:"size:2000"`                        // 当绑定自己的邮件或手机号码时, 表示为密码
	RefreshToken string `json:"refresh_token" gorm:"size:2000"`
	SocialName   string `json:"social_name"`
	Gender       string `json:"gender"`
	IconUrl      string `json:"icon_url"`
}

func GetBindCacheKey(bindType, bindId string) string {
	return strings.Join([]string{"binding", bindType, bindId}, ":")
}
