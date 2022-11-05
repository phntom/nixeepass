package webui

import (
	"github.com/mileusna/useragent"
	"github.com/phntom/nixeepass/orm"
	"github.com/rs/zerolog/log"
	"net/http"
	"regexp"
	"strings"
)

var webdavMatcher = regexp.MustCompile("^/webdav/([0-9a-z-]{36})/my.kdbx$")

func webdavHandler(writer http.ResponseWriter, request *http.Request) {
	match := webdavMatcher.FindStringSubmatch(request.URL.Path)
	if len(match) != 2 {
		enrich(log.Info(), request).Msg("invalid request url")
		return
	}
	token := match[1]
	db := orm.GetDB()
	var device orm.Device
	result := db.First(&device, orm.Device{Token: token})
	if result.Error != nil {
		enrich(log.Info().Err(result.Error).Str("token", token), request).Msg("token not found")
		return
	}
	ua := useragent.Parse(request.Header.Get("User-Agent"))
	if ua.Name == "okhttp" {
		ua.OS = "android"
		ua.Name = "keepass2android"
		if ua.Version == "4.10.0-RC1" {
			ua.Version = "1.09c"
		}
	}
	device.UserAgentDevice = strings.ToLower(ua.Name)
	device.UserAgentOS = strings.ToLower(ua.OS)
	device.UserAgentOSVersion = ua.OSVersion
	device.UserAgentDeviceVersion = ua.Version
	device.LastIP = request.RemoteAddr[:strings.LastIndexByte(request.RemoteAddr, ':')]
	device.LastCountry = "israel"
	result = db.Updates(&device)
	if result.Error != nil {
		enrich(log.Error().Err(result.Error).Int("device_id", int(device.ID)).Str("user_id", device.UserID), request).Msg("failed to update device record")
		return
	}
	enrich(log.Info().Err(result.Error).Int("device_id", int(device.ID)).Str("user_id", device.UserID).Str("os", device.UserAgentOS).Str("ua", device.UserAgentDevice).Str("ip", device.LastIP).Str("country", device.LastCountry), request).Msg("database downloaded")
}
