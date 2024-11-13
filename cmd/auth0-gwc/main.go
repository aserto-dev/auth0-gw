package main

import (
	"context"
	"os"

	"github.com/auth0/go-auth0/management"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
)

func main() {
	var (
		eventType string
		target    string
		source    string
		userEmail string
	)
	flag.StringVarP(
		&eventType, "event-type", "e", "",
		"event type [post-login | post-change-password | pre-user-registration | post-user-registration | sync]",
	)
	flag.StringVarP(&userEmail, "user-email", "u", "", "user email")
	flag.StringVarP(&target, "target", "t", "", "target url")
	flag.StringVarP(&source, "source", "s", "auth0.com", "event source")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if target == "" {
		log.Fatal().Msg("--target not set")
	}
	if source == "" {
		log.Fatal().Msg("--source not set")
	}
	if eventType == "" {
		log.Fatal().Msg("--event-type not set")
	}

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatal().Err(err).Msg("create http client")
	}

	user := management.User{
		Email: &userEmail,
	}

	// Create an Event.
	event := cloudevents.NewEvent()
	event.SetSource(source)
	event.SetType(eventType)
	if err := event.SetData(cloudevents.ApplicationJSON, user); err != nil {
		log.Error().Err(err).Msg("set data")
		os.Exit(1)
	}

	// Set a target.
	ctx := cloudevents.ContextWithTarget(context.Background(), target)

	// Send that Event.
	if result := c.Send(ctx, event); cloudevents.IsUndelivered(result) {
		log.Error().Msgf("send %v", result)
	} else {
		log.Info().Msgf("sent: %v", event)
		log.Info().Msgf("result: %v", result)
	}
}
