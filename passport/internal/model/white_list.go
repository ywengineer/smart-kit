package model

import "gitee.com/ywengineer/smart-kit/pkg/rdbs"

type WhiteList struct {
	rdbs.Model
	Passport uint `json:"passport" redis:"passport" gorm:"size:50;unique;comment:账号ID"`
}
