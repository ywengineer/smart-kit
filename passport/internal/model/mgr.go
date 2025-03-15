package model

import "strconv"

type MgrUser struct {
	Model
	Account  string `json:"account" redis:"account" gorm:"index:idx_account;not null;comment:账号"`
	Password string `json:"password" redis:"password" gorm:"not null;comment:密码"`
	Name     string `json:"name" redis:"name" gorm:"comment:姓名;size:20"`
	DeptNo   uint   `json:"dept_no" redis:"dept_no" gorm:"comment:部门;default:0"`
	Title    string `json:"title" redis:"title" gorm:"comment:职位"`
}

func GetMgrCacheKey(account string) string {
	return "mgr:" + account
}

func GetWhiteListCacheKey(id int64) string {
	return "white-list:" + strconv.FormatInt(id, 10)
}
