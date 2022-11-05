package webui

import (
	"bytes"
	_ "embed"
	"github.com/dustin/go-humanize"
	"github.com/fsnotify/fsnotify"
	"github.com/lpar/gzipped/v2"
	"github.com/phntom/nixeepass/config"
	"github.com/rs/zerolog/log"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
)

//go:embed root.html
var rootContent string
var rootTemplate *template.Template

type httpConfig struct {
	BindAddress       string `mapstructure:"listen"`
	DashboardEndpoint string `mapstructure:"dashboard_endpoint"`
}

var cfg *httpConfig

func RunWebUI() error {
	err := loadTemplates()
	if err != nil {
		return err
	}
	fs := http.StripPrefix("/static/", gzipped.FileServer(gzipped.Dir("./static")))
	http.Handle("/static/", fs)
	http.HandleFunc("/", rootHandler)
	cfg = &httpConfig{}
	ReloadConfig(fsnotify.Event{})
	config.RegisterOnConfigChange(ReloadConfig)

	log.Info().Str("BindAddress", cfg.BindAddress).Msg("Starting web server")
	return http.ListenAndServe(cfg.BindAddress, nil)
}

func loadTemplates() error {
	rootTemplate = template.New("").Funcs(template.FuncMap{
		"ToUpper":       strings.ToUpper,
		"ToLower":       strings.ToLower,
		"Title":         GetTitleMap,
		"FA":            GetIconMap,
		"HumanizeTime":  HumanTime,
		"HumanizeBytes": humanize.Bytes,
	})
	_, err := rootTemplate.Parse(rootContent)
	return err
}

func ReloadConfig(in fsnotify.Event) {
	log.Debug().Msg("Reloading http config")
	err := config.Config().UnmarshalKey("http", cfg)
	if err != nil {
		log.Info().Err(err).Msg("failed to reload http config")
	}
	if cfg.BindAddress == "" {
		cfg.BindAddress = ":8979"
	}
	if cfg.DashboardEndpoint == "" {
		cfg.DashboardEndpoint = "/dashboard"
	}
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	var buf bytes.Buffer
	var err error
	if request.URL.Path == cfg.DashboardEndpoint {
		err = dashboardHandler(&buf, request)
	} else if request.URL.Path == "/" {
		http.Redirect(writer, request, cfg.DashboardEndpoint, http.StatusFound)
		return
	} else {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Info().Err(err).Msg("failed to render dashboard")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/html; charset=UTF-8")
	_, err = buf.WriteTo(writer)
	if err != nil {
		log.Info().Err(err).Msg("failed to transmit dashboard buffer")
	}
}

func dashboardHandler(writer io.Writer, request *http.Request) error {
	never := time.Unix(0, 0)
	now := time.Now().Add(-10 * time.Minute)
	yesterday := time.Now().Add(-24 * time.Hour)
	lastMonth := time.Now().Add(-24 * 40 * time.Hour)
	lastYear := time.Now().Add(-24 * 365 * time.Hour)

	data := RootPage{
		NeedsLogin: false,
		AppName:    "nixeepass",
		BrandName:  "kix.co.il",
		UserID:     "abc123",
		Devices: []DeviceDetails{
			{
				Created:  now,
				Modified: never,
				Browser:  "chrome",
				Country:  "united states",
				OS:       "windows",
			},
			{
				Created:  yesterday,
				Modified: now,
				Browser:  "firefox",
				Country:  "canada",
				OS:       "ubuntu",
			},
			{
				Created:  lastYear,
				Modified: lastMonth,
				Browser:  "keepass2android",
				Country:  "unknown",
				OS:       "android",
			},
			{
				Created:  lastMonth,
				Modified: now,
				Browser:  "opera",
				Country:  "israel",
				OS:       "macos",
			},
		},
		BackupActiveModified: now,
		Backups: []BackupDetails{
			{
				BackupID: "123",
				Modified: now,
				IsActive: true,
			},
			{
				BackupID: "456",
				Modified: yesterday,
				IsActive: false,
			},
		},
	}
	return rootTemplate.Execute(writer, data)
}
