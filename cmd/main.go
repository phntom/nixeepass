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
	logJSON := cfg.GetBool("log.json")
	logDebug := cfg.GetBool("log.debug")
	logTrace := cfg.GetBool("log.trace")

	level := zerolog.InfoLevel
	jwwLevel := jww.LevelInfo
	if logDebug {
		level = zerolog.DebugLevel
		jwwLevel = jww.LevelDebug
	}
	if logTrace {
		level = zerolog.TraceLevel
		jwwLevel = jww.LevelTrace
	}

	zerolog.SetGlobalLevel(level)

	if !logJSON {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
		jww.SetStdoutOutput(consoleWriter)
		jww.SetStdoutThreshold(jwwLevel)
		jww.SetPrefix("viper")
		log.Logger = log.Output(consoleWriter)
	}
	zerolog.DefaultContextLogger = &log.Logger

	log.WithLevel(level).Msg(fmt.Sprintf("Starting %s", config.AppName))
}
