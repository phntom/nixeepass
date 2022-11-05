package webui

import (
	"github.com/phntom/nixeepass/orm"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

func revokePerformer(writer http.ResponseWriter, request *http.Request) error {
	id, err := strconv.ParseUint(request.PostFormValue("id"), 10, 32)
	if err != nil {
		return err
	}

	url := request.URL
	url.Path = cfgHttp.DashboardEndpoint
	http.Redirect(writer, request, url.String(), http.StatusFound)

	var device orm.Device
	db := orm.GetDB()
	tx := db.First(&device, orm.Device{ID: uint(id), UserID: "abc123"})
	if tx.Error != nil {
		enrich(log.Info().Err(err), request).Str("user_id", device.UserID).Msg("revoke record not found")
		return nil
	}
	db.Delete(&device)
	if tx.Error != nil {
		enrich(log.Info().Err(err), request).Str("user_id", device.UserID).Msg("failed to update record for revoke")
		return nil
	}
	enrich(log.Info().Int("key", int(device.ID)).Str("user_id", device.UserID), request).Msg("revoked device")

	return nil
}
