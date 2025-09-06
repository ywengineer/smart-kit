package model

import (
	"time"

	"gitee.com/ywengineer/smart-kit/pkg/rdbs"
)

type BaseLog struct {
	rdbs.Model       `json:",inline"`
	GameId           string    `json:"game_id" redis:"game_id"`                       // 游戏ID
	ServerId         string    `json:"server_id" redis:"server_id"`                   // 服务器标识
	Passport         string    `json:"passport" redis:"passport"`                     // 游戏账号
	PlayerId         string    `json:"player_id" redis:"player_id"`                   // 角色ID
	LogicId          string    `json:"logic_id" redis:"logic_id"`                     // 逻辑ID
	PlayerName       string    `json:"player_name" redis:"player_name"`               // 角色名称
	Extra            string    `json:"extra" redis:"extra" gorm:"type:text"`          // 额外数据
	SystemType       string    `json:"system_type" redis:"system_type"`               // 账号注册的系统类型
	Locale           string    `json:"locale" redis:"locale"`                         // 账号注册地区
	PlayerCreateTime time.Time `json:"player_create_time" redis:"player_create_time"` // 角色注册时间
}
