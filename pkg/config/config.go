package config

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Gateway   *Gateway   `json:"gateway"`   // gateway settings
	Directory *Directory `json:"directory"` // target directory settings
	Auth0     *Auth0     `json:"auth0"`     // source auth0 settings
	Loader    *Loader    `json:"loader"`    // ds-load* settings
	Scheduler *Scheduler `json:"scheduler"` // scheduler settings
}

type Gateway struct {
	Port int    `json:"port"` // gateway listen port
	Path string `json:"path"` // gateway URL path
}

type Directory struct {
	Host     string `json:"host"`      // directory gRPC host address
	APIKey   string `json:"api_key"`   // directory read-write API key
	TenantID string `json:"tenant_id"` // directory tenant ID
	Insecure bool   `json:"insecure"`  // skip TLS validation
}

type Auth0 struct {
	Domain       string `json:"domain"`        // Auth0 domain
	ClientID     string `json:"client_id"`     // Auth0 client ID
	ClientSecret string `json:"client_secret"` // Auth0 client secret
	Connection   string `json:"connection"`    // Auth0 connection name
	UserPID      string `json:"user_pid"`      // Auth0 user ID of user to export
	UserEmail    string `json:"user_email"`    // Auth0 email name of user to export
	InclRoles    bool   `json:"incl_roles"`    // Auth0 include roles in export
}

type Loader struct {
	BinPath  string `json:"bin_path"` // ds-load absolute path to binaries (ds-load & ds-load-auth0)
	Template string `json:"template"` // ds-load absolute path to template file
}

type Scheduler struct {
	Interval string `json:"interval"` // time interval string as 5m30s for 5 min and 30 seconds, minimum interval is 1m.
}

func Load(file string) (*Config, error) {
	v := viper.New()

	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetConfigFile(file)

	v.SetDefault("scheduler.interval", "5m")

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.Wrapf(err, "loading config file %s", file)
	}

	cfg := &Config{}
	err := v.UnmarshalExact(cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "json"
	})
	if err != nil {
		return nil, errors.Wrapf(err, "loading config file %s", file)
	}

	if dur, err := time.ParseDuration(cfg.Scheduler.Interval); err != nil || dur < (time.Minute*1) {
		log.Warn().Str("new interval", "1m").Str("old interval", cfg.Scheduler.Interval).Msg("reset, min is 1m")
		cfg.Scheduler.Interval = "1m"
	}

	return cfg, nil
}
