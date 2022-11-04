package webui

import (
	_ "embed"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"net/http"
	"strings"
)

//go:embed root.html
var rootContent string
var rootTemplate *template.Template

func RunWebUI() error {
	err := loadTemplates()
	if err != nil {
		return err
	}
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	http.Handle("/static/", fs)
	http.HandleFunc("/", rootHandler)
	return http.ListenAndServe(":8999", nil)
}

func loadTemplates() error {
	rootTemplate = template.New("").Funcs(template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
		"Title":   cases.Title(language.English).String,
		"FA":      GetIconMap,
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
	data := RootPage{
		//LoginSection:   "is-active",
		//ContentSection: "is-hidden",
		LoginSection:   "is-hidden",
		ContentSection: "is-active",
		BrandName:      "kix.co.il",
		UserID:         "abc123",
		Devices: []DeviceDetails{
			{
				Browser: "chrome",
				Country: "united states",
				OS:      "windows",
			},
			{
				Browser: "firefox",
				Country: "canada",
				OS:      "ubuntu",
			},
			{
				Browser: "keepass2android",
				Country: "unknown",
				OS:      "android",
			},
			{
				Browser: "opera",
				Country: "israel",
				OS:      "mac",
			},
		},
	}
	_ = rootTemplate.Execute(writer, data)
}
