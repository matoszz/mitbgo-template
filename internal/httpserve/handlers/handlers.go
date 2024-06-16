package handlers

import (
	echo "github.com/datumforge/echox"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/datumforge/datum/pkg/providers/github"
	"github.com/datumforge/datum/pkg/providers/google"
	"github.com/datumforge/datum/pkg/providers/webauthn"
	"github.com/datumforge/datum/pkg/sessions"
	"github.com/lestrrat-go/jwx/v2/jwk"

	ent "github.com/datumforge/go-template/internal/ent/generated"
)

// Handler contains configuration options for handlers
type Handler struct {
	// IsTest is a flag to determine if the application is running in test mode and will mock external calls
	IsTest bool
	// DBClient to interact with the generated ent schema
	DBClient *ent.Client
	// RedisClient to interact with redis
	RedisClient *redis.Client
	// Logger provides the zap logger to do logging things from the handlers
	Logger *zap.SugaredLogger
	// ReadyChecks is a set of checkFuncs to determine if the application is "ready" upon startup
	ReadyChecks Checks
	// SessionConfig to handle sessions
	SessionConfig *sessions.SessionConfig
	// AuthMiddleware contains the middleware to be used for authenticated endpoints
	AuthMiddleware []echo.MiddlewareFunc
	// JWTKeys contains the set of valid JWT authentication key
	JWTKeys jwk.Set
	// OauthProvider contains the configuration settings for all supported Oauth2 providers
	OauthProvider OauthProviderConfig
}

// OauthProviderConfig represents the configuration for OAuth providers such as Github and Google
type OauthProviderConfig struct {
	// RedirectURL is the URL that the OAuth2 client will redirect to after authentication with datum
	RedirectURL string `json:"redirectUrl" koanf:"redirectUrl" default:"http://localhost:3001/api/auth/callback/datum"`
	// Github contains the configuration settings for the Github Oauth Provider
	Github github.ProviderConfig `json:"github" koanf:"github"`
	// Google contains the configuration settings for the Google Oauth Provider
	Google google.ProviderConfig `json:"google" koanf:"google"`
	// Webauthn contains the configuration settings for the Webauthn Oauth Provider
	Webauthn webauthn.ProviderConfig `json:"webauthn" koanf:"webauthn"`
}
