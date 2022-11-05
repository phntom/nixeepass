package webui

import (
	"bytes"
	"github.com/dustin/go-humanize"
	"github.com/fsnotify/fsnotify"
	"github.com/lpar/gzipped/v2"
	"github.com/phntom/nixeepass/config"
	"github.com/phntom/nixeepass/orm"
	"github.com/phntom/nixeepass/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
	"strings"
)

var cfgHttp *httpConfig
var cfgDashboard *dashboardConfig

func RunWebUI() error {

	orm.StartORM()

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
	dashboardTemplate = template.New("").Funcs(template.FuncMap{
		"ToUpper":       strings.ToUpper,
		"ToLower":       strings.ToLower,
		"Title":         utils.GetTitleMap,
		"FA":            GetIconMap,
		"HumanizeTime":  utils.HumanTime,
		"HumanizeBytes": humanize.Bytes,
	})
	_, err := dashboardTemplate.Parse(dashboardContent)
	return err
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	var buf bytes.Buffer
	var err error
	if request.URL.Path == "/" {
		http.Redirect(writer, request, cfgHttp.DashboardEndpoint, http.StatusFound)
		enrich(log.Debug(), request).Msg("redirecting to dashboard")
		return
	} else if request.URL.Path == cfgHttp.DashboardEndpoint {
		err = dashboardHandler(&buf, request)
	} else if request.URL.Path == "/_liveliness" {
		err = livelinessHandler(&buf, request)
	} else if request.URL.Path == "/_readiness" {
		err = readinessHandler(&buf, request)
	} else {
		enrich(log.Debug(), request).Msg("redirecting to dashboard")
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		enrich(log.Error().Err(err), request).Msg("failed to render dashboard")
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

func enrich(l *zerolog.Event, request *http.Request) *zerolog.Event {
	if cfgHttp.Log.Method && len(request.Method) > 0 {
		l = l.Str("method", request.Method)
	}
	if cfgHttp.Log.RemoteAddr && len(request.RemoteAddr) > 0 {
		l = l.Str("remote_addr", request.RemoteAddr)
	}
	for _, header := range cfgHttp.Log.Headers {
		h := request.Header.Get(header)
		if len(h) > 0 {
			l = l.Str(header, h)
		}
	}
	if len(request.URL.Host) > 0 {
		l = l.Str("host", request.URL.Host)
	}
	return l.Str("path", request.URL.Path)
}
