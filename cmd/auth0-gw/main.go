package main

import (
	"context"
	"os"

	"github.com/aserto-dev/auth0-gw/pkg/config"
	"github.com/auth0/go-auth0/management"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/magefile/mage/sh"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfg *config.Config

func main() {
	var configFile string
	flag.StringVarP(&configFile, "config", "c", "", "path to config.yaml file")
	flag.Parse()

	v := viper.New()

	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetConfigFile(configFile)

	if err := v.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msgf("loading config file %s", configFile)
	}

	cfg = new(config.Config)
	err := v.UnmarshalExact(cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "json"
	})
	if err != nil {
		log.Fatal().Err(err).Msgf("loading config file %s", configFile)
	}

	ctx := context.Background()
	p, err := cloudevents.NewHTTP(http.WithPort(cfg.Gateway.Port), http.WithPath(cfg.Gateway.Path))
	if err != nil {
		log.Error().Err(err).Msg("failed to create protocol")
	}

	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Error().Err(err).Msg("failed to create client")
	}

	log.Info().Msgf("listen on :%d%s\n", p.Port, p.Path)
	if err := c.StartReceiver(ctx, receive); err != nil {
		log.Fatal().Err(err).Msg("failed to start listener")
	}
}

func receive(ctx context.Context, event cloudevents.Event) {
	switch event.Context.GetType() {
	case "sync":
		log.Printf("sync:\n")
		go sync(cfg, "")

	case "post-login":
		var u management.User
		if err := u.UnmarshalJSON(event.Data()); err != nil {
			log.Printf("err: %s\n", err)
		}
		log.Printf("post-login:\n")
		log.Printf("user: %s\n", u.GetEmail())

		go sync(cfg, *u.Email)
		// log.Printf("%s\n", event)

	case "post-change-password":
		var u management.User
		if err := u.UnmarshalJSON(event.Data()); err != nil {
			log.Printf("err: %s\n", err)
		}
		log.Printf("post-change-password:\n")
		log.Printf("user: %s\n", u.GetEmail())

		go sync(cfg, *u.Email)
		// log.Printf("%s\n", event)

	case "pre-user-registration":
		var u management.User
		if err := u.UnmarshalJSON(event.Data()); err != nil {
			log.Printf("err: %s\n", err)
		}
		log.Printf("pre-user-registration:\n")
		log.Printf("user: %s\n", u.GetEmail())

		go sync(cfg, *u.Email)
		// log.Printf("%s\n", event)

	case "post-user-registration":
		var u management.User
		if err := u.UnmarshalJSON(event.Data()); err != nil {
			log.Printf("err: %s\n", err)
		}
		log.Printf("post-user-registration:\n")
		log.Printf("user: %s\n", u.GetEmail())

		go sync(cfg, *u.Email)
		// log.Printf("%s\n", event)

	default:
		log.Printf("default handler: %s", event)
	}
}

func sync(cfg *config.Config, email string) error {
	env := map[string]string{
		"DS_TEMPLATE_FILE":      cfg.Loader.Template,
		"DIRECTORY_HOST":        cfg.Directory.Host,
		"DIRECTORY_API_KEY":     cfg.Directory.APIKey,
		"DIRECTORY_TENANT_ID":   cfg.Directory.TenantID,
		"AUTH0_DOMAIN":          cfg.Auth0.Domain,
		"AUTH0_CLIENT_ID":       cfg.Auth0.ClientID,
		"AUTH0_CLIENT_SECRET":   cfg.Auth0.ClientSecret,
		"AUTH0_CONNECTION_NAME": cfg.Auth0.Connection,
		"AUTH0_USER_PID":        cfg.Auth0.UserPID,
		"AUTH0_USER_EMAIL":      email,
		"AUTH0_ROLES":           "false",
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cmd := "ds-load"
	args := []string{"exec", "auth0"}

	ran, err := sh.Exec(env, os.Stdout, os.Stderr, cmd, args...)
	if !ran {
		log.Warn().Msgf("command %s did not run", cmd)
	}
	if err != nil {
		log.Err(err).Msgf("command %s failed", cmd)
	}
	return nil
}
