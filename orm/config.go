package orm

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/phntom/nixeepass/config"
	"github.com/rs/zerolog/log"
)

type ormConfig struct {
	Type    string           `mapstructure:"type"`
	Sqlite3 ormConfigSqlite3 `mapstructure:"sqlite3"`
}

type ormConfigSqlite3 struct {
	Filename          string `mapstructure:"filename"`
	PreparedStatement bool   `mapstructure:"prepared_statement"`
}

var cfg ormConfig

func reloadConfig(in fsnotify.Event) {
	log.Debug().Msg("Reloading orm configs")
	err := config.Config().UnmarshalKey("db", &cfg)
	if err != nil {
		log.Info().Err(err).Msg("failed to reload db config")
	}
	if cfg.Type == "" {
		cfg.Type = "sqlite3"
	}
	if cfg.Sqlite3.Filename == "" {
		cfg.Sqlite3.Filename = fmt.Sprintf("%s.db", config.AppName)
	}
	reconnect()
}
