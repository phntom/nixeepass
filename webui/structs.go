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
	Devices              []DeviceDetails
	BackupActiveModified time.Time
	Backups              []BackupDetails
}
