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
	"time"
)

var cfgHttp *httpConfig
var cfgDashboard *dashboardConfig

func RunWebUI() error {

	orm.StartORM()

	err := loadIP2Location()
	if err != nil {
		return err
	}

	err = loadTemplates()
	if err != nil {
		return err
	}
	fs := http.StripPrefix("/static/", gzipped.FileServer(gzipped.Dir("./static")))
	http.Handle("/static/", fs)
	http.HandleFunc("/webdav/", webdavHandler)
	http.HandleFunc("/", rootHandler)
	cfgHttp = &httpConfig{}
	cfgDashboard = &dashboardConfig{}
	reloadConfig(fsnotify.Event{})
	config.RegisterOnConfigChange(reloadConfig)

	log.Info().Str("BindAddress", cfgHttp.BindAddress).Msg("Starting web server")
	return http.ListenAndServe(cfgHttp.BindAddress, nil)
}

func loadTemplates() error {
	funcMap := template.FuncMap{
		"ToUpper":       strings.ToUpper,
		"ToLower":       strings.ToLower,
		"Title":         utils.GetTitleMap,
		"FA":            GetIconMap,
		"HumanizeTime":  utils.HumanTime,
		"HumanizeBytes": humanize.Bytes,
		"timeEQ":        time.Time.Equal,
	}
	dashboardTemplate = template.New("dashboard").Funcs(funcMap)
	_, err := dashboardTemplate.Parse(dashboardContent)
	if err != nil {
		return err
	}
	grantTemplate = template.New("grant").Funcs(funcMap)
	_, err = grantTemplate.Parse(grantContent)
	return err
}

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	var buf bytes.Buffer
	var err error
	if request.URL.Path == "/" && request.Method == http.MethodGet {
		http.Redirect(writer, request, cfgHttp.Endpoints.Dashboard, http.StatusFound)
		enrich(log.Debug(), request).Msg("redirecting to dashboard")
		return
	} else if request.URL.Path == cfgHttp.Endpoints.Dashboard && request.Method == http.MethodGet {
		err = dashboardHandler(&buf, request)
	} else if request.URL.Path == cfgHttp.Endpoints.Grant && request.Method == http.MethodPost {
		err = grantPerformer(writer, request)
	} else if request.URL.Path == cfgHttp.Endpoints.Grant && request.Method == http.MethodGet {
		err = grantHandler(&buf, request)
	} else if request.URL.Path == cfgHttp.Endpoints.Revoke && request.Method == http.MethodPost {
		err = revokePerformer(writer, request)
	} else if request.URL.Path == cfgHttp.Endpoints.Login && request.Method == http.MethodPost {
	} else if request.URL.Path == cfgHttp.Endpoints.Login && request.Method == http.MethodGet {
		http.Redirect(writer, request, cfgHttp.Endpoints.Dashboard, http.StatusFound)
		enrich(log.Info(), request).Msg("http get login attempt, redirecting to dashboard")
		return
	} else if request.URL.Path == cfgHttp.Endpoints.Liveliness && request.Method == http.MethodGet {
		err = livelinessHandler(&buf, request)
	} else if request.URL.Path == cfgHttp.Endpoints.Readiness && request.Method == http.MethodGet {
		err = readinessHandler(&buf, request)
	} else {
		enrich(log.Debug(), request).Msg("not found")
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		enrich(log.Error().Err(err), request).Msg("failed to render dashboard")
		http.Error(writer, template.HTMLEscapeString(err.Error()), http.StatusInternalServerError)
		return
	}
	if buf.Len() == 0 {
		return
	}
	writer.WriteHeader(http.StatusOK)
	for key, val := range cfgHttp.Headers {
		writer.Header().Set(key, val)
	}
	_, err = buf.WriteTo(writer)
	if err != nil {
		enrich(log.Info().Err(err), request).Msg("failed to transmit dashboard buffer")
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
