package main

import (
	"fmt"
	"github.com/phntom/nixeepass/config"
	"github.com/phntom/nixeepass/webui"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	jww "github.com/spf13/jwalterweatherman"
	"os"
)

func main() {
	StartLogging()
	log.Err(webui.RunWebUI())
}

func StartLogging() {
	cfg := config.Config()
	logJson := cfg.GetBool("log.json")
	logDebug := cfg.GetBool("log.debug")

	level := zerolog.InfoLevel
	jwwLevel := jww.LevelInfo
	if logDebug {
		level = zerolog.DebugLevel
		jwwLevel = jww.LevelDebug
	}

	zerolog.SetGlobalLevel(level)

	if !logJson {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
		jww.SetLogOutput(consoleWriter)
		jww.SetStdoutThreshold(jwwLevel)
		jww.SetLogThreshold(jwwLevel)
		jww.SetPrefix("viper")
		log.Logger = log.Output(consoleWriter)
	}

	log.WithLevel(level).Msg(fmt.Sprintf("Starting %s", config.AppName))
}
