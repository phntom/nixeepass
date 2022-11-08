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
	BindAddress string `mapstructure:"listen"`
	Endpoints   httpConfigEndpoints
	Headers     map[string]string
	Log         httpConfigLogs
	Auth        httpConfigAuth
}

type httpConfigEndpoints struct {
	Dashboard                 string
	Grant                     string
	Revoke                    string
	Login                     string
	Liveliness                string
	Readiness                 string
	Webdav                    string
	LivelinessReadinessSecret string
	PublicPrefix              string
}

type httpConfigLogs struct {
	Method     bool
	RemoteAddr bool
	Headers    []string
}

type httpConfigAuth struct {
	MockUserID string `mapstructure:"mock_user_id"`
}

type dashboardConfig struct {
	AppName   string
	BrandName string
	IconMap   map[string]string `mapstructure:"icons"`
}
