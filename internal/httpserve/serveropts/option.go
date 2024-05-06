package serveropts

import (
	echoprometheus "github.com/datumforge/echo-prometheus/v5"
	echo "github.com/datumforge/echox"
	"github.com/datumforge/echox/middleware"
	"github.com/datumforge/echozap"
	"github.com/datumforge/entx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/datumforge/go-template/internal/ent/generated"
	"github.com/datumforge/go-template/internal/graphapi"
	"github.com/datumforge/go-template/internal/httpserve/config"
	"github.com/datumforge/go-template/internal/httpserve/server"

	"github.com/datumforge/datum/pkg/cache"
	authmw "github.com/datumforge/datum/pkg/middleware/auth"
	"github.com/datumforge/datum/pkg/middleware/cachecontrol"
	"github.com/datumforge/datum/pkg/middleware/cors"
	"github.com/datumforge/datum/pkg/middleware/echocontext"
	"github.com/datumforge/datum/pkg/middleware/mime"
	"github.com/datumforge/datum/pkg/middleware/ratelimit"
	"github.com/datumforge/datum/pkg/middleware/redirect"
	"github.com/datumforge/datum/pkg/middleware/sentry"
	"github.com/datumforge/datum/pkg/sessions"
)

type ServerOption interface {
	apply(*ServerOptions)
}

type applyFunc struct {
	applyInternal func(*ServerOptions)
}

func (fso *applyFunc) apply(s *ServerOptions) {
	fso.applyInternal(s)
}

func newApplyFunc(apply func(option *ServerOptions)) *applyFunc {
	return &applyFunc{
		applyInternal: apply,
	}
}

// WithConfigProvider supplies the config for the server
func WithConfigProvider(cfgProvider config.ConfigProvider) ServerOption {
	return newApplyFunc(func(s *ServerOptions) {
		s.ConfigProvider = cfgProvider
	})
}

// WithLogger supplies the logger for the server
func WithLogger(l *zap.SugaredLogger) ServerOption {
	return newApplyFunc(func(s *ServerOptions) {
		// Add logger to main config
		s.Config.Logger = l
		// Add logger to the handlers config
		s.Config.Handler.Logger = l
	})
}

// WithHTTPS sets up TLS config settings for the server
func WithHTTPS() ServerOption {
	return newApplyFunc(func(s *ServerOptions) {
		if !s.Config.Settings.Server.TLS.Enabled {
			// this is set to enabled by WithServer
			// if TLS is not enabled, move on
			return
		}

		s.Config.WithTLSDefaults()

		if !s.Config.Settings.Server.TLS.AutoCert {
			s.Config.WithTLSCerts(s.Config.Settings.Server.TLS.CertFile, s.Config.Settings.Server.TLS.CertKey)
		}
	})
}

// WithReadyChecks adds readiness checks to the server
func WithReadyChecks(c *entx.EntClientConfig, r *redis.Client) ServerOption {
	return newApplyFunc(func(s *ServerOptions) {
		// Always add a check to the primary db connection
		s.Config.Handler.AddReadinessCheck("db_primary", entx.Healthcheck(c.GetPrimaryDB()))

		// Check the secondary db, if enabled
		if s.Config.Settings.DB.MultiWrite {
			s.Config.Handler.AddReadinessCheck("db_secondary", entx.Healthcheck(c.GetSecondaryDB()))
		}

		// Check the connection to redis, if enabled
		if s.Config.Settings.Redis.Enabled {
			s.Config.Handler.AddReadinessCheck("redis", cache.Healthcheck(r))
		}
	})
}

// WithGraphRoute adds the graph handler to the server
func WithGraphRoute(srv *server.Server, c *generated.Client) ServerOption {
	return newApplyFunc(func(s *ServerOptions) {
		// Setup Graph API Handlers
		r := graphapi.NewResolver(c).
			WithLogger(s.Config.Logger.Named("resolvers"))

		handler := r.Handler(s.Config.Settings.Server.Dev)

		// Add Graph Handler
		srv.AddHandler(handler)
	})
}

// WithMiddleware adds the middleware to the server
func WithMiddleware() ServerOption {
	return newApplyFunc(func(s *ServerOptions) {
		// Initialize middleware if null
		if s.Config.DefaultMiddleware == nil {
			s.Config.DefaultMiddleware = []echo.MiddlewareFunc{}
		}

		redirectMW := redirect.Config{
			Redirects: map[string]string{
				"/.well-known/change-password": "/v1/forgot-password",
				"/security.txt":                "/.well-known/security.txt",
			},
			Code: 302, // nolint: gomnd
		}
		// default middleware
		s.Config.DefaultMiddleware = append(s.Config.DefaultMiddleware,
			middleware.RequestID(), // add request id
			middleware.Recover(),   // recover server from any panic/fatal error gracefully
			middleware.LoggerWithConfig(middleware.LoggerConfig{
				Format: "remote_ip=${remote_ip}, method=${method}, uri=${uri}, status=${status}, session=${header:Set-Cookie}, host=${host}, referer=${referer}, user_agent=${user_agent}, route=${route}, path=${path}, auth=${header:Authorization}\n",
			}),
			sentry.New(),
			echoprometheus.MetricsMiddleware(),                   // add prometheus metrics
			echozap.ZapLogger(s.Config.Logger.Desugar()),         // add zap logger, middleware requires the "regular" zap logger
			echocontext.EchoContextToContextMiddleware(),         // adds echo context to parent
			cors.New(s.Config.Settings.Server.CORS.AllowOrigins), // add cors middleware
			mime.NewWithConfig(mime.Config{DefaultContentType: echo.MIMEApplicationJSONCharsetUTF8}), // add mime middleware
			cachecontrol.New(),                 // add cache control middleware
			middleware.Secure(),                // add XSS middleware
			redirect.NewWithConfig(redirectMW), // add redirect middleware
		)
	})
}

// WithRateLimiter sets up the rate limiter for the server
func WithRateLimiter() ServerOption {
	return newApplyFunc(func(s *ServerOptions) {
		if s.Config.Settings.Ratelimit.Enabled {
			s.Config.DefaultMiddleware = append(s.Config.DefaultMiddleware, ratelimit.RateLimiterWithConfig(&s.Config.Settings.Ratelimit))
		}
	})
}

// WithSessionManager sets up the default session manager with a 10 minute ttl
// with persistence to redis
func WithSessionManager(rc *redis.Client) ServerOption {
	return newApplyFunc(func(s *ServerOptions) {
		cc := sessions.DefaultCookieConfig

		// In order for things to work in dev mode with localhost
		// we need to se the debug cookie config
		if s.Config.Settings.Server.Dev {
			cc = &sessions.DebugOnlyCookieConfig
		} else {
			cc.Name = sessions.DefaultCookieName
		}

		sm := sessions.NewCookieStore[map[string]any](cc,
			[]byte(s.Config.Settings.Sessions.SigningKey),
			[]byte(s.Config.Settings.Sessions.EncryptionKey),
		)

		// add session middleware, this has to be added after the authMiddleware so we have the user id
		// when we get to the session. this is also added here so its only added to the graph routes
		// REST routes are expected to add the session middleware, as required
		sessionConfig := sessions.NewSessionConfig(
			sm,
			sessions.WithPersistence(rc),
			sessions.WithLogger(s.Config.Logger),
			sessions.WithSkipperFunc(authmw.SessionSkipperFunc),
		)

		// set cookie config to be used
		sessionConfig.CookieConfig = cc

		// Make the cookie session store available
		// to graph and REST endpoints
		s.Config.Handler.SessionConfig = &sessionConfig
		s.Config.SessionConfig = &sessionConfig

		s.Config.GraphMiddleware = append(s.Config.GraphMiddleware,
			sessions.LoadAndSaveWithConfig(sessionConfig),
		)
	})
}
