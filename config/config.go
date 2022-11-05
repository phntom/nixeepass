package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
	"time"

	"github.com/spf13/viper"
)

const AppName = "nixeepass"

var onConfigChange []func(fsnotify.Event)

// Provider defines a set of read-only methods for accessing the application
// configuration params as defined in one of the config files.
type Provider interface {
	ConfigFileUsed() string
	Get(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
	Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) error
	UnmarshalKey(key string, rawVal interface{}, opts ...viper.DecoderConfigOption) error
}

var defaultConfig *viper.Viper

// Config returns a default config providers
func Config() Provider {
	return defaultConfig
}

// LoadConfigProvider returns a configured viper instance
func LoadConfigProvider(appName string) Provider {
	return readViperConfig(appName)
}

func init() {
	defaultConfig = readViperConfig(AppName)
}

func readViperConfig(appName string) *viper.Viper {
	v := viper.New()
	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	v.SetConfigName(fmt.Sprintf("%s-config", appName))

	// docker mount points
	v.AddConfigPath(fmt.Sprintf("/etc/%s/", appName))
	v.AddConfigPath(fmt.Sprintf("$HOME/.%s", appName))

	// for development
	v.AddConfigPath("config")

	err := v.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Panic().Err(err).Msg("fatal error config file")
	}
	v.OnConfigChange(OnConfigChange)
	v.WatchConfig()

	return v
}

func OnConfigChange(e fsnotify.Event) {
	log.Info().Str("filename", e.Name).Int("len(onConfigChange)", len(onConfigChange)).Msg("Config file changed")
	for _, f := range onConfigChange {
		f(e)
	}
}

func RegisterOnConfigChange(run func(in fsnotify.Event)) {
	onConfigChange = append(onConfigChange, run)
}
