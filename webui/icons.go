package webui

var iconMap = map[string]string{
	"chrome":          "fa-brands fa-chrome",
	"firefox":         "fa-brands fa-firefox",
	"opera":           "fa-brands fa-opera",
	"android":         "fa-brands fa-android",
	"keepass2android": "fa-solid fa-mobile",
	"mac":             "fa-brands fa-apple",
	"ubuntu":          "fa-brands fa-ubuntu",
	"windows":         "fa-brands fa-windows",
	"canada":          "fa-brands fa-canadian-maple-leaf",
	"israel":          "fa-solid fa-star-of-david",
	"united states":   "fa-solid fa-flag-usa",
	"unknown":         "fa-solid fa-globe",
}

func GetIconMap(key string) string {
	if val, ok := iconMap[key]; ok {
		return val
	}
	return "fa-regular fa-circle-question"
}
