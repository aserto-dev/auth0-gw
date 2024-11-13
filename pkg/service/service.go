package service

import (
	"context"

	"github.com/aserto-dev/auth0-gw/pkg/config"
	"github.com/aserto-dev/auth0-gw/pkg/loader"

	"github.com/auth0/go-auth0/management"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Service struct {
	cfg  *config.Config
	load *loader.Loader
}

func New(cfg *config.Config) *Service {
	return &Service{
		cfg:  cfg,
		load: loader.New(cfg),
	}
}

func (s *Service) Start(ctx context.Context) error {
	p, err := cloudevents.NewHTTP(
		http.WithPort(s.cfg.Gateway.Port),
		http.WithPath(s.cfg.Gateway.Path),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create protocol")
	}

	c, err := cloudevents.NewClient(p)
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	log.Info().Msgf("listen on :%d%s", p.Port, p.Path)
	if err := c.StartReceiver(ctx, s.dispatch); err != nil {
		return errors.Wrap(err, "failed to start listener")
	}

	return nil
}

func (s *Service) Stop() {
}

func (s *Service) dispatch(ctx context.Context, event cloudevents.Event) {
	log.Info().Str("type", event.Context.GetType()).Msg("dispatch")

	var u management.User
	switch event.Context.GetType() {
	case "post-login", "post-change-password", "pre-user-registration", "post-user-registration":
		if err := u.UnmarshalJSON(event.Data()); err != nil {
			log.Error().Err(err).Msg("unmarshal event data")
			return
		}

		go func() {
			req := *s.cfg.Auth0
			if u.ID != nil && *u.ID != "" {
				req.UserPID = *u.ID
			} else if u.Email != nil && *u.Email != "" {
				req.UserEmail = *u.Email
			} else {
				log.Error().Msg("event user section does not contain ID or email")
				return
			}
			s.load.Sync(&req)
		}()

	case "sync":
		go func() {
			req := *s.cfg.Auth0
			req.UserPID = ""
			req.UserEmail = ""
			s.load.Sync(&req)
		}()

	default:
		log.Error().Msgf("unknown event type %s", event.Context.GetType())
		return
	}
}
