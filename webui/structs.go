package webui

import (
	"github.com/phntom/nixeepass/orm"
	"time"
)

type DeviceDetails struct {
	Country  string
	Browser  string
	OS       string
	Created  time.Time
	Modified time.Time
}

type BackupDetails struct {
	BackupID string
	Modified time.Time
	IsActive bool
}

type RootPage struct {
	NeedsLogin           bool
	AppName              string
	BrandName            string
	UserID               string
	UserName             string
	Devices              []orm.Device
	BackupActiveModified time.Time
	Backups              []orm.Backup
}

type httpConfig struct {
	BindAddress       string            `mapstructure:"listen"`
	DashboardEndpoint string            `mapstructure:"dashboard_endpoint"`
	Headers           map[string]string `mapstructure:"headers"`
	Log               httpConfigLogs    `mapstructure:"log"`
}

type httpConfigLogs struct {
	Method     bool     `mapstructure:"method"`
	RemoteAddr bool     `mapstructure:"remote_addr"`
	Headers    []string `mapstructure:"headers"`
}

type dashboardConfig struct {
	AppName   string            `mapstructure:"app_name"`
	BrandName string            `mapstructure:"brand_name"`
	IconMap   map[string]string `mapstructure:"icons"`
}
