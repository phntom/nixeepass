package webui

import (
	_ "embed" // it's ok
	"fmt"
	"github.com/google/uuid"
	"github.com/phntom/nixeepass/orm"
	"github.com/rs/zerolog/log"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

//go:embed grant.html
var grantContent string
var grantTemplate *template.Template

func grantPerformer(writer http.ResponseWriter, request *http.Request) error {
	db := orm.GetDB()
	device := orm.Device{
		UserID: "abc123",
		Token:  uuid.New().String(),
	}
	tx := db.Create(&device)
	if tx.Error != nil {
		enrich(log.Error().Err(tx.Error).Str("user_id", device.UserID), request).Msg("grant failed to create record")
		return tx.Error
	}
	enrich(log.Info().Int("key", int(device.ID)).Str("user_id", device.UserID), request).Msg("granted new device")

	url := request.URL
	q := url.Query()
	q.Set("id", fmt.Sprintf("%d", device.ID))
	url.RawQuery = q.Encode()
	http.Redirect(writer, request, url.String(), http.StatusFound)
	return nil
}

func grantHandler(writer io.Writer, request *http.Request) error {
	id, err := strconv.ParseUint(request.URL.Query().Get("id"), 10, 32)
	if err != nil {
		return err
	}
	var device orm.Device
	tx := orm.GetDB().First(&device, orm.Device{ID: uint(id), UserID: "abc123"})
	if tx.Error != nil {
		return tx.Error
	}
	data := GrantPage{
		AppName:    cfgDashboard.AppName,
		BrandName:  cfgDashboard.BrandName,
		CSRF:       "abc123",
		HTTPConfig: cfgHttp,
		UserID:     "abc123",
		UserName:   "phantom@kix.co.il",
		NeedsLogin: false,
		Device:     &device,
		URL:        request.URL,
	}
	return grantTemplate.Execute(writer, data)
}
