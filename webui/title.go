package webui

import (
	"github.com/dustin/go-humanize"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"time"
)

var titleMap = map[string]string{
	"keepass2android": "Keepass2Android",
	"macos":           "macOS",
}

var neverTime = time.Unix(0, 0)

func GetTitleMap(key string) string {
	if val, ok := titleMap[key]; ok {
		return val
	}
	return cases.Title(language.English).String(key)
}

func HumanTime(then time.Time) string {
	if then == neverTime {
		return "Never"
	}
	human := humanize.Time(then)
	return strings.ToUpper(human[0:1]) + human[1:]
}
