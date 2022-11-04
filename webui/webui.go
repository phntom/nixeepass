package webui

import (
	_ "embed"
	"github.com/dustin/go-humanize"
	"github.com/lpar/gzipped/v2"
	"html/template"
	"net/http"
	"strings"
	"time"
)

//go:embed root.html
var rootContent string
var rootTemplate *template.Template

func RunWebUI() error {
	err := loadTemplates()
	if err != nil {
		return err
	}
	fs := http.StripPrefix("/static/", gzipped.FileServer(gzipped.Dir("./static")))
	http.Handle("/static/", fs)
	http.HandleFunc("/", rootHandler)
	return http.ListenAndServe(":8999", nil)
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

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)

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
	_ = rootTemplate.Execute(writer, data)
}
