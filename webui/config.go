package webui

import (
	"github.com/fsnotify/fsnotify"
	"github.com/phntom/nixeepass/config"
	"github.com/rs/zerolog/log"
)

func reloadConfig(in fsnotify.Event) {
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
	if cfgHttp.DashboardEndpoint == "" {
		cfgHttp.DashboardEndpoint = "/dashboard"
	}
	if cfgDashboard.AppName == "" {
		cfgDashboard.AppName = "nixeepass"
	}
	if cfgDashboard.BrandName == "" {
		cfgDashboard.BrandName = "kix.co.il"
	}
}
