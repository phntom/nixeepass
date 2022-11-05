package orm

import (
	"github.com/fsnotify/fsnotify"
	"github.com/phntom/nixeepass/config"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"reflect"
)

var db *gorm.DB

func StartORM() {
	log.Debug().Msg("Starting ORM")
	reloadConfig(fsnotify.Event{})
	config.RegisterOnConfigChange(reloadConfig)
}

func register() {
	var registerEntities = []interface{}{
		Device{},
		Backup{},
	}
	for _, entity := range registerEntities {
		err := db.AutoMigrate(&entity)
		if err != nil {
			log.Error().Err(err).Str("entity", reflect.TypeOf(entity).Name()).Msg("Failed to migrate orm")
		}
	}
}

func Alive() error {
	var device Device
	result := db.First(&device)
	if result.Error.Error() == "record not found" {
		return nil
	}
	return result.Error
}

func GetDB() *gorm.DB {
	return db
}
