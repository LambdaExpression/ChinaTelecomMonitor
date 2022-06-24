package models

import (
	"github.com/golang-module/carbon/v2"
)

type Summary struct {
	ID           int             `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string          `gorm:"type:varchar(1024)" json:"username"` // 电信账号名
	Use          int64           `gorm:"type:int" json:"use"`                // 流量使用量
	Total        int64           `gorm:"type:int" json:"total"`              // 流量总量
	GeneralUse   int64           `gorm:"type:int" json:"generalUse"`         // 通用流量使用量
	GeneralTotal int64           `gorm:"type:int" json:"generalTotal"`       // 通用流量总量
	SpecialUse   int64           `gorm:"type:int" json:"specialUse"`         // 专用流量使用量
	SpecialTotal int64           `gorm:"type:int" json:"specialTotal"`       // 专用流量总量
	Balance      int64           `gorm:"type:int" json:"balance"`            // 余额
	VoiceUsage   int64           `gorm:"type:int" json:"voiceUsage"`         // 语音使用量
	VoiceAmount  int64           `gorm:"type:int" json:"voiceAmount"`        // 语音总量
	CreateTime   carbon.DateTime `gorm:"type:TIMESTAMP" json:"createTime"`   // 查询时间
	Items        []SummaryItems  `json:"items"`
}

type SummaryItems struct {
	Name  string `gorm:"type:varchar(1024)" json:"name"`
	Use   int64  `gorm:"type:int" json:"use"`
	Total int64  `gorm:"type:int" json:"total"`
}
