package scheduler

import (
	"context"
	"time"

	"github.com/aserto-dev/auth0-gw/pkg/config"
	"github.com/aserto-dev/auth0-gw/pkg/loader"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

const tickInterval = time.Second * 15

type Scheduler struct {
	cfg  *config.Config
	load *loader.Loader
}

func New(cfg *config.Config) *Scheduler {
	sched := gocron.NewScheduler(time.UTC)
	sched.SingletonModeAll()

	return &Scheduler{
		cfg:  cfg,
		load: loader.New(cfg),
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	log.Info().Msg("start scheduler")

	go s.next(ctx, time.NewTicker(tickInterval))

	return nil
}

func (s *Scheduler) Stop() {
}

func (s *Scheduler) next(ctx context.Context, interval *time.Ticker) {
	defer interval.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Warn().Time("done", time.Now().UTC()).Msg("scheduler")
			return
		case t := <-interval.C:
			log.Info().Time("dispatch", t).Time("now", time.Now().UTC()).Msg("scheduler")

			s.task()

			interval.Stop()

			nextInterval, err := time.ParseDuration(s.cfg.Scheduler.Interval)
			if err != nil {
				log.Error().Err(err).Msg("parse next interval duration")
				return
			}

			interval = time.NewTicker(nextInterval)

			log.Info().Str("interval", s.cfg.Scheduler.Interval).Time("next", time.Now().UTC().Add(nextInterval)).Msg("scheduler")
		}
	}
}

func (s *Scheduler) task() {
	req := *s.cfg.Auth0
	req.UserEmail = ""
	req.UserPID = ""
	s.load.Sync(&req)
}
