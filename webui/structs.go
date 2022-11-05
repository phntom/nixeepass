package webui

import "time"

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
	Devices              []DeviceDetails
	BackupActiveModified time.Time
	Backups              []BackupDetails
}

type httpConfig struct {
	BindAddress       string            `mapstructure:"listen"`
	DashboardEndpoint string            `mapstructure:"dashboard_endpoint"`
	Headers           map[string]string `mapstructure:"headers"`
}

type dashboardConfig struct {
	AppName   string            `mapstructure:"app_name"`
	BrandName string            `mapstructure:"brand_name"`
	IconMap   map[string]string `mapstructure:"icons"`
}
