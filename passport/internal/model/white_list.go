package model

type WhiteList struct {
	Model
	Passport uint `json:"passport" redis:"passport" gorm:"size:50;unique;comment:账号ID"`
}
