package orm

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func reconnect() {
	var err error
	log.Debug().Str("type", cfg.Type).Msg("Connecting to ORM")
	if cfg.Type == "sqlite3" {
		log.Info().Str("type", cfg.Type).Str("filename", cfg.Sqlite3.Filename).Msg("Opening sqlite3 file")
		db, err = gorm.Open(sqlite.Open(cfg.Sqlite3.Filename), &gorm.Config{
			PrepareStmt: cfg.Sqlite3.PreparedStatement,
			Logger:      Logger{},
		})
		if err != nil {
			log.Error().Err(err).Str("filename", cfg.Sqlite3.Filename).Msg("Failed to open sqlite3 db")
		}
	} else {
		log.Error().Str("type", cfg.Type).Msg("Unknown DB type")
	}
	register()
}
