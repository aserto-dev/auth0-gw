package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/aserto-dev/auth0-gw/pkg/config"
	"github.com/aserto-dev/auth0-gw/pkg/scheduler"
	"github.com/aserto-dev/auth0-gw/pkg/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
)

func main() {
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("starting auth0-gw")

	var configFile string
	flag.StringVarP(&configFile, "config", "c", "", "path to config.yaml file")
	flag.Parse()

	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("loading config")
	}

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	go func() {
		sched := scheduler.New(ctx, cfg)
		defer sched.Stop()

		if err := sched.Start(); err != nil {
			log.Fatal().Err(err).Msg("starting scheduler")
		}
	}()

	go func() {
		svc := service.New(ctx, cfg)
		defer svc.Stop()

		if err := svc.Start(); err != nil {
			log.Fatal().Err(err).Msg("starting service")
		}
	}()

	// wait for context to be canceled
	<-ctx.Done()
}
