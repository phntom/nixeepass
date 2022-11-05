package webui

import (
	"io"
	"net/http"
	"time"
)

func dashboardHandler(writer io.Writer, request *http.Request) error {
	never := time.Unix(0, 0)
	now := time.Now().Add(-10 * time.Minute)
	yesterday := time.Now().Add(-24 * time.Hour)
	lastMonth := time.Now().Add(-24 * 40 * time.Hour)
	lastYear := time.Now().Add(-24 * 365 * time.Hour)

	data := RootPage{
		NeedsLogin: false,
		AppName:    cfgDashboard.AppName,
		BrandName:  cfgDashboard.BrandName,
		UserID:     "abc123",
		UserName:   "NameHere",
		Devices: []DeviceDetails{
			{
				Created:  now,
				Modified: never,
				Browser:  "chrome",
				Country:  "united states",
				OS:       "windows",
			},
			{
				Created:  yesterday,
				Modified: now,
				Browser:  "firefox",
				Country:  "canada",
				OS:       "ubuntu",
			},
			{
				Created:  lastYear,
				Modified: lastMonth,
				Browser:  "keepass2android",
				Country:  "unknown",
				OS:       "android",
			},
			{
				Created:  lastMonth,
				Modified: now,
				Browser:  "opera",
				Country:  "israel",
				OS:       "macos",
			},
		},
		BackupActiveModified: now,
		Backups: []BackupDetails{
			{
				BackupID: "123",
				Modified: now,
				IsActive: true,
			},
			{
				BackupID: "456",
				Modified: yesterday,
				IsActive: false,
			},
		},
	}
	return rootTemplate.Execute(writer, data)
}

func GetIconMap(key string) string {
	if val, ok := cfgDashboard.IconMap[key]; ok {
		return val
	}
	return "fa-regular fa-circle-question"
}
