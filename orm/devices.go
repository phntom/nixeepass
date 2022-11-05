package orm

import (
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	Active          bool   `gorm:"index:active_user_id,priority:1"`
	UserID          string `gorm:"index:active_user_id,priority:2"`
	Token           string `gorm:"uniqueIndex:token"`
	UserAgentOS     string `json:"user_agent_os"`
	UserAgentDevice string `json:"user_agent_device"`
	LastIP          string `json:"last_ip"`
	LastCountry     string `json:"last_country"`
}
