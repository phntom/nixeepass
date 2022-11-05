package orm

import (
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	ID                     uint   `gorm:"primarykey"`
	UserID                 string `gorm:"index:device_user_id"`
	Token                  string `gorm:"uniqueIndex:token"`
	UserAgentOS            string `json:"user_agent_os"`
	UserAgentOSVersion     string `json:"user_agent_os_version"`
	UserAgentDevice        string `json:"user_agent_device"`
	UserAgentDeviceVersion string `json:"user_agent_device_version"`
	LastIP                 string `json:"last_ip"`
	LastCountry            string `json:"last_country"`
}
