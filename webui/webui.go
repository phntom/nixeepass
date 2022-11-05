package webui

import (
	"bytes"
	_ "embed"
	"github.com/dustin/go-humanize"
	"github.com/fsnotify/fsnotify"
	"github.com/lpar/gzipped/v2"
	"github.com/phntom/nixeepass/config"
	"github.com/phntom/nixeepass/utils"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
	"strings"
)

//go:embed root.html
var rootContent string
var rootTemplate *template.Template

var cfgHttp *httpConfig
var cfgDashboard *dashboardConfig

func RunWebUI() error {
	err := loadTemplates()
	if err != nil {
		return err
	}
	fs := http.StripPrefix("/static/", gzipped.FileServer(gzipped.Dir("./static")))
	http.Handle("/static/", fs)
	http.HandleFunc("/", rootHandler)
	cfgHttp = &httpConfig{}
	cfgDashboard = &dashboardConfig{}
	reloadConfig(fsnotify.Event{})
	config.RegisterOnConfigChange(reloadConfig)

	log.Info().Str("BindAddress", cfgHttp.BindAddress).Msg("Starting web server")
	return http.ListenAndServe(cfgHttp.BindAddress, nil)
}

func loadTemplates() error {
	rootTemplate = template.New("").Funcs(template.FuncMap{
		"ToUpper":       strings.ToUpper,
		"ToLower":       strings.ToLower,
		"Title":         utils.GetTitleMap,
		"FA":            GetIconMap,
		"HumanizeTime":  utils.HumanTime,
		"HumanizeBytes": humanize.Bytes,
	})
	_, err := rootTemplate.Parse(rootContent)
	return err
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	var buf bytes.Buffer
	var err error
	if request.URL.Path == cfgHttp.DashboardEndpoint {
		err = dashboardHandler(&buf, request)
	} else if request.URL.Path == "/" {
		http.Redirect(writer, request, cfgHttp.DashboardEndpoint, http.StatusFound)
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
	for key, val := range cfgHttp.Headers {
		writer.Header().Set(key, val)
	}
	_, err = buf.WriteTo(writer)
	if err != nil {
		log.Info().Err(err).Msg("failed to transmit dashboard buffer")
	}
}
