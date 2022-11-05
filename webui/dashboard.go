package webui

import (
	_ "embed"
	"github.com/phntom/nixeepass/orm"
	"html/template"
	"io"
	"net/http"
	"time"
)

//go:embed dashboard.html
var dashboardContent string
var dashboardTemplate *template.Template

func dashboardHandler(writer io.Writer, request *http.Request) error {
	userID := "abc123"
	userName := "phantom@kix.co.il"

	var devices []orm.Device
	db := orm.GetDB()
	result := db.Find(&devices, orm.Device{UserID: userID, Active: true}, 100)
	err := result.Error
	if err != nil {
		return err
	}
	var backups []orm.Backup
	result = db.Find(&backups, orm.Backup{UserID: userID}, 20)
	err = result.Error
	if err != nil {
		return err
	}
	backupActiveModified := time.Unix(0, 0)
	if len(backups) > 0 {
		backupActiveModified = backups[0].CreatedAt
	}

	data := RootPage{
		NeedsLogin:           false,
		AppName:              cfgDashboard.AppName,
		BrandName:            cfgDashboard.BrandName,
		UserID:               userID,
		UserName:             userName,
		Devices:              devices,
		BackupActiveModified: backupActiveModified,
		Backups:              backups,
	}
	return dashboardTemplate.Execute(writer, data)
}

func GetIconMap(key string) string {
	if val, ok := cfgDashboard.IconMap[key]; ok {
		return val
	}
	return "fa-regular fa-circle-question"
}
