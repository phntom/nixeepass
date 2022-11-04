package webui

type DeviceDetails struct {
	Country     string
	CountryFlag string
	Browser     string
	BrowserIcon string
	OS          string
	OSIcon      string
}

type RootPage struct {
	BrandName      string
	LoginSection   string
	ContentSection string
	UserID         string
	Devices        []DeviceDetails
}
