package loader

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aserto-dev/auth0-gw/pkg/config"
	"github.com/magefile/mage/sh"
	"github.com/rs/zerolog/log"
)

type Loader struct {
	config *config.Config
}

func New(cfg *config.Config) *Loader {
	return &Loader{
		config: cfg,
	}
}

func (l *Loader) Sync(a *config.Auth0) {
	log.Info().Str("email", a.UserEmail).Str("pid", a.UserPID).Msg("sync started")

	env := map[string]string{
		"DS_TEMPLATE_FILE":      l.config.Loader.Template,
		"DIRECTORY_HOST":        l.config.Directory.Host,
		"DIRECTORY_API_KEY":     l.config.Directory.APIKey,
		"DIRECTORY_TENANT_ID":   l.config.Directory.TenantID,
		"AUTH0_DOMAIN":          a.Domain,
		"AUTH0_CLIENT_ID":       a.ClientID,
		"AUTH0_CLIENT_SECRET":   a.ClientSecret,
		"AUTH0_CONNECTION_NAME": a.Connection,
		"AUTH0_USER_PID":        a.UserPID,
		"AUTH0_USER_EMAIL":      a.UserEmail,
		"AUTH0_ROLES":           strconv.FormatBool(a.InclRoles),
	}

	// create fully qualified name to ds-load binary.
	cmd := filepath.Join(l.config.Loader.BinPath, "ds-load")

	// --no-rate-limit => use the auth0 SDK ratelimiting implementation.
	args := []string{"exec", "auth0", "--no-rate-limit"}

	// track duration and log result when finished.
	defer l.timer(time.Now(), a, "sync finished")

	// spawn ds-load to execute fetch, transform and import into directory.
	ran, err := sh.Exec(env, os.Stdout, os.Stderr, cmd, args...)
	if !ran {
		log.Warn().Msgf("command %s did not run", cmd)
	}
	if err != nil {
		log.Err(err).Msgf("command %s failed", cmd)
	}
}

func (l *Loader) timer(start time.Time, a *config.Auth0, msg string) {
	log.Info().Str("email", a.UserEmail).Str("pid", a.UserPID).TimeDiff("duration", time.Now(), start).Msg(msg)
}
