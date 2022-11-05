package webui

import (
	_ "embed"
	"github.com/ip2location/ip2location-go/v9"
	"github.com/mileusna/useragent"
	"github.com/phntom/nixeepass/orm"
	"github.com/rs/zerolog/log"
	"net/http"
	"regexp"
	"strings"
)

var webdavMatcher = regexp.MustCompile("^/webdav/([0-9a-z-]{36})/my.kdbx$")
var ip2locationDB *ip2location.DB

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

	parseUA(&device, request)
	parseIP(&device, request)
	err := parseCountry(&device)
	if err != nil {
		enrich(log.Error().Err(result.Error).Int("device_id", int(device.ID)).Str("user_id", device.UserID), request).Msg("failed to resolve ip to country")
	}

	result = db.Updates(&device)
	if result.Error != nil {
		enrich(log.Error().Err(result.Error).Int("device_id", int(device.ID)).Str("user_id", device.UserID), request).Msg("failed to update device record")
		return
	}
	enrich(log.Info().Err(result.Error).Int("device_id", int(device.ID)).Str("user_id", device.UserID).Str("os", device.UserAgentOS).Str("ua", device.UserAgentDevice).Str("ip", device.LastIP).Str("country", device.LastCountry), request).Msg("database downloaded")
}

func parseCountry(device *orm.Device) error {
	ip := device.LastIP
	record, err := ip2locationDB.Get_country_long(ip)
	if err != nil {
		return err
	}
	device.LastCountry = strings.ToLower(record.Country_long)
	if device.LastCountry == "-" {
		device.LastCountry = "global"
	}
	return nil
}

func parseIP(device *orm.Device, request *http.Request) {
	device.LastIP = request.RemoteAddr[:strings.LastIndexByte(request.RemoteAddr, ':')]
}

func parseUA(device *orm.Device, request *http.Request) {
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
}

func loadIP2Location() error {
	var err error

	ip2locationDB, err = ip2location.OpenDB("IP2LOCATION-LITE-DB1.IPV6.BIN")
	if err != nil {
		log.Error().Err(err).Msg("Failed loading ip2location database")
	} else {
		log.Info().Msg("Loaded ip2location database")
	}
	return err
}
