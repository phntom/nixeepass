package webui

import (
	"github.com/phntom/nixeepass/orm"
	"net/url"
	"time"
)

type DashboardPage struct {
	CSRF                 string
	NeedsLogin           bool
	AppName              string
	BrandName            string
	UserID               string
	UserName             string
	Devices              []orm.Device
	BackupActiveModified time.Time
	Backups              []orm.Backup
	HTTPConfig           *httpConfig
}

type GrantPage struct {
	CSRF       string
	HTTPConfig *httpConfig
	AppName    string
	BrandName  string
	NeedsLogin bool
	UserID     string
	UserName   string
	Device     *orm.Device
	URL        *url.URL
}

type httpConfig struct {
	BindAddress               string            `mapstructure:"listen"`
	DashboardEndpoint         string            `mapstructure:"dashboard_endpoint"`
	GrantEndpoint             string            `mapstructure:"grant_endpoint"`
	RevokeEndpoint            string            `mapstructure:"revoke_endpoint"`
	LoginEndpoint             string            `mapstructure:"login_endpoint"`
	LivelinessEndpoint        string            `mapstructure:"liveliness_endpoint"`
	ReadinessEndpoint         string            `mapstructure:"readiness_endpoint"`
	WebDavEndpoint            string            `mapstructure:"webdav_endpoint"`
	LivelinessReadinessSecret string            `mapstructure:"liveliness_readiness_secret"`
	PublicPrefix              string            `mapstructure:"public_prefix"`
	Headers                   map[string]string `mapstructure:"headers"`
	Log                       httpConfigLogs    `mapstructure:"log"`
}

type httpConfigLogs struct {
	Method     bool     `mapstructure:"method"`
	RemoteAddr bool     `mapstructure:"remote_addr"`
	Headers    []string `mapstructure:"headers"`
}

type dashboardConfig struct {
	AppName   string            `mapstructure:"app_name"`
	BrandName string            `mapstructure:"brand_name"`
	IconMap   map[string]string `mapstructure:"icons"`
}
