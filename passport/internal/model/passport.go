package model

import (
	"gitee.com/ywengineer/smart-kit/pkg/rdbs"
	"strconv"
	"strings"
	"time"
)
import "gorm.io/datatypes"

type Passport struct {
	rdbs.Model
	DeviceId   string         `json:"device_id" redis:"device_id" gorm:"index:idx_device;comment:设备ID"`
	Adid       string         `json:"adid" redis:"adid" gorm:"comment:设备广告标识"`
	SystemType string         `json:"system_type" redis:"system_type" gorm:"comment:系统类型;size:20"`
	Locale     string         `json:"locale" redis:"locale" gorm:"comment:地区"`
	Extra      datatypes.JSON `json:"extra" redis:"extra" gorm:"comment:主机额外信息"`
}

type PassportPunish struct {
	rdbs.Model
	PassportId uint      `json:"passport_id" gorm:"index:idx_passport;comment:账号ID" redis:"passport_id"`
	DeviceId   string    `json:"device_id" gorm:"index:idx_device;comment:设备ID" redis:"device_id"`
	Type       string    `json:"type" gorm:"comment:惩罚类型" redis:"type"`
	BeginTime  time.Time `json:"begin_time"  redis:"begin_time"`
	EndTime    time.Time `json:"end_time" redis:"end_time"`
	Reason     string    `json:"reason" gorm:"size:200;comment:原因" redis:"reason"`
}

func GetPassportCacheKey(id uint) string {
	return strings.Join([]string{"passport", strconv.FormatUint(uint64(id), 10)}, ":")
}
