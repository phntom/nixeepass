package webui

import (
	"github.com/fsnotify/fsnotify"
	"github.com/phntom/nixeepass/config"
	"github.com/rs/zerolog/log"
)

func reloadConfig(_ fsnotify.Event) {
	log.Debug().Msg("Reloading http, dashboard configs")
	err := config.Config().UnmarshalKey("http", cfgHttp)
	if err != nil {
		log.Info().Err(err).Msg("failed to reload http config")
	}
	err = config.Config().UnmarshalKey("dashboard", cfgDashboard)
	if err != nil {
		log.Info().Err(err).Msg("failed to reload dashboard config")
	}
	if cfgHttp.BindAddress == "" {
		cfgHttp.BindAddress = ":8979"
	}
	if cfgHttp.Endpoints.Dashboard == "" {
		cfgHttp.Endpoints.Dashboard = "/dashboard"
	}
	if cfgHttp.Endpoints.Grant == "" {
		cfgHttp.Endpoints.Grant = "/grant"
	}
	if cfgHttp.Endpoints.Revoke == "" {
		cfgHttp.Endpoints.Revoke = "/revoke"
	}
	if cfgHttp.Endpoints.Login == "" {
		cfgHttp.Endpoints.Login = "/login"
	}
	if cfgHttp.Endpoints.Liveliness == "" {
		cfgHttp.Endpoints.Liveliness = "/_liveliness"
	}
	if cfgHttp.Endpoints.Readiness == "" {
		cfgHttp.Endpoints.Readiness = "/_readiness"
	}
	if cfgHttp.Endpoints.Webdav == "" {
		cfgHttp.Endpoints.Webdav = "/webdav"
	}
	if cfgHttp.Endpoints.PublicPrefix == "" {
		cfgHttp.Endpoints.PublicPrefix = "https://example.com"
	}
	if cfgDashboard.AppName == "" {
		cfgDashboard.AppName = "nixeepass"
	}
	if cfgDashboard.BrandName == "" {
		cfgDashboard.BrandName = "kix.co.il"
	}
}
