package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/aserto-dev/auth0-gw/pkg/config"
	"github.com/aserto-dev/auth0-gw/pkg/scheduler"
	"github.com/aserto-dev/auth0-gw/pkg/service"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
)

func main() {
	var configFile string
	var templateFile string
	var consoleMode bool
	flag.StringVarP(&configFile, "config", "c", "", "path to config.yaml file")
	flag.StringVarP(&templateFile, "template", "t", "", "path to transform.tmpl file")
	flag.BoolVarP(&consoleMode, "console", "", false, "console mode")
	flag.Parse()

	if consoleMode {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	log.Info().Msg("starting auth0-gw")

	if ok, err := fileExists(configFile); !ok {
		log.Fatal().Err(err).Msgf("config file %s not found", configFile)
	}

	if ok, err := fileExists(templateFile); !ok {
		log.Fatal().Err(err).Msgf("template file %s not found", templateFile)
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("loading config")
	}
	cfg.Loader.Template = templateFile

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	go func() {
		sched := scheduler.New(cfg)
		defer sched.Stop()

		if err := sched.Start(ctx); err != nil {
			log.Fatal().Err(err).Msg("starting scheduler")
		}
	}()

	go func() {
		svc := service.New(cfg)
		defer svc.Stop()

		if err := svc.Start(ctx); err != nil {
			log.Fatal().Err(err).Msg("starting service")
		}
	}()

	// wait for context to be canceled
	<-ctx.Done()
}

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, errors.Wrapf(err, "failed to stat file '%s'", path)
	}
}
