package model

import "gorm.io/gorm"

type WhiteList struct {
	gorm.Model
	Passport uint `json:"passport" gorm:"size:50;unique;comment:账号ID"`
}
