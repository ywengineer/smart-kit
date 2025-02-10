package model

import (
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)
import "gorm.io/datatypes"

type Passport struct {
	gorm.Model
	DeviceId   string         `json:"device_id" gorm:"index:idx_device;comment:设备ID"`
	Adid       string         `json:"adid" gorm:"index:idx_adid;comment:设备广告标识"`
	SystemType string         `json:"system_type" gorm:"comment:系统类型;size:20"`
	Locale     string         `json:"locale" gorm:"comment:地区"`
	Extra      datatypes.JSON `json:"extra" gorm:"comment:主机额外信息"`
}

type PassportPunish struct {
	gorm.Model
	PassportId uint      `json:"passport_id" gorm:"index:idx_passport;comment:账号ID"`
	DeviceId   string    `json:"device_id" gorm:"index:idx_device;comment:设备ID"`
	Type       string    `json:"type" gorm:"comment:惩罚类型"`
	BeginTime  time.Time `json:"begin_time" `
	EndTime    time.Time `json:"end_time" `
	Reason     string    `json:"reason" gorm:"size:200;comment:原因"`
}

func GetPassportCacheKey(id uint) string {
	return strings.Join([]string{"passport", strconv.FormatUint(uint64(id), 10)}, ":")
}
