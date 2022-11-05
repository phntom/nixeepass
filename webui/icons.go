package webui

var iconMap = map[string]string{
	"chrome":          "fa-brands fa-chrome",
	"firefox":         "fa-brands fa-firefox",
	"opera":           "fa-brands fa-opera",
	"android":         "fa-brands fa-android",
	"keepass2android": "fa-solid fa-mobile",
	"macos":           "fa-brands fa-apple",
	"ubuntu":          "fa-brands fa-ubuntu",
	"windows":         "fa-brands fa-windows",
	"canada":          "fa-brands fa-canadian-maple-leaf",
	"israel":          "fa-solid fa-star-of-david",
	"united states":   "fa-solid fa-flag-usa",
	"united kingdom":  "fa-solid fa-sterling-sign",
	"japan":           "fa-solid fa-yen-sign",
	"brazil":          "fa-solid fa-brazilian-real-sign",
	"france":          "fa-solid fa-archway",
	"georgia":         "fa-solid fa-lari-sign",
	"indonesian":      "fa-solid fa-rupiah-sign",
	"india":           "fa-solid fa-indian-rupee-sign",
	"russia":          "fa-solid fa-ruble-sign",
	"ukraine":         "fa-solid fa-hryvnia-sign",
	"mexico":          "fa-solid fa-peso-sign",
	"nigeria":         "fa-solid fa-naira-sign",
	"kazakhstan":      "fa-solid fa-tenge-sign",
	"azerbaijan":      "fa-solid fa-manat-sign",
	"turkey":          "fa-solid fa-turkish-lira-sign",
	"laos":            "fa-solid fa-kip-sign",
	"paraguay":        "fa-solid fa-guarani-sign",
	"vietnam":         "fa-solid fa-dong-sign",
	"costa rica":      "fa-solid fa-colon-sign",
	"ghana":           "fa-solid fa-cedi-sign",
	"thailand":        "fa-solid fa-baht-sign",
	"switzerland":     "fa-solid fa-franc-sign",
	"south korean":    "fa-solid fa-won-sign",
	"euro":            "fa-solid fa-euro-sign",
	"global":          "fa-solid fa-globe",
	"oceania":         "fa-solid fa-earth-oceania",
	"europe":          "fa-solid fa-earth-europe",
	"asia":            "fa-solid fa-earth-asia",
	"africa":          "fa-solid fa-earth-africa",
	"americas":        "fa-solid fa-earth-americas",
	"unknown":         "fa-solid fa-circle-question",
	"space":           "fa-solid fa-satellite",
	"air":             "fa-solid fa-plane-up",
	"ocean":           "fa-solid fa-sailboat",
}

func GetIconMap(key string) string {
	if val, ok := iconMap[key]; ok {
		return val
	}
	return "fa-regular fa-circle-question"
}
