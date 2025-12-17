package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/oriiyx/fritz/app/core/services/session"
	services "github.com/oriiyx/fritz/app/core/services/user"
	"github.com/oriiyx/fritz/app/core/utils/env"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/rs/zerolog"
)

type API struct {
	DB             *pgxpool.Pool
	Goth           *goth.Session
	Store          *sessions.CookieStore
	Logger         *zerolog.Logger
	EnvConf        *env.Conf
	Queries        *db.Queries
	Validator      *validator.Validate
	UserService    *services.UserService
	SessionManager *session.Manager
}

// Define a custom type for context keys to avoid collisions
type contextKey string

const (
	ProviderKey      = "provider"
	SessionID        = "session.id"
	UserIDKey        = "user.id"
	AuthenticatedKey = "authenticated"

	// Context key using the custom type
	providerContextKey contextKey = "provider"
)

func NewAuth(logger *zerolog.Logger, validator *validator.Validate, cfg *env.Conf, pool *pgxpool.Pool, store *sessions.CookieStore, queries *db.Queries) *API {
	gothic.Store = store
	goth.UseProviders(
		google.New(cfg.GoogleAuth.ID, cfg.GoogleAuth.Secret, fmt.Sprintf("%s/auth/google/callback", cfg.GetBaseURL()), "email", "profile"),
		github.New(cfg.GithubAuth.ID, cfg.GithubAuth.Secret, fmt.Sprintf("%s/auth/github/callback", cfg.GetBaseURL()), "read:user", "user:email"),
	)

	return &API{
		DB:             pool,
		Store:          store,
		Logger:         logger,
		EnvConf:        cfg,
		Queries:        db.New(pool),
		Validator:      validator,
		UserService:    services.NewUserService(pool),
		SessionManager: session.NewManager(cfg, queries, logger),
	}
}

func (a *API) ProviderLogin(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, ProviderKey)
	r = r.WithContext(context.WithValue(r.Context(), providerContextKey, provider))

	session, _ := a.Store.Get(r, SessionID)
	isAuthenticated := session.Values[AuthenticatedKey]
	// oAuthUser := session.Values[UserIDKey]
	if isAuthenticated != nil && isAuthenticated == true {
		http.Redirect(w, r, a.EnvConf.GetBaseURL(), http.StatusFound)
	}

	// try to get the user without re-authenticating
	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		a.Logger.Info().Interface("user", user).Msg("Got user information")
		dbUser, err := a.Queries.GetUserByEmail(r.Context(), user.Email)
		if err != nil {
			a.Logger.Err(err).Str("email", user.Email).Msg("Could not find user by email")
			http.Redirect(w, r, a.EnvConf.GetBaseURL(), http.StatusInternalServerError)
			return
		}
		session.Values[AuthenticatedKey] = true
		session.Values[UserIDKey] = dbUser.ID.String()
		err = session.Save(r, w)
		if err != nil {
			a.Logger.Err(err).Msg("Could not save session")
			http.Redirect(w, r, a.EnvConf.GetBaseURL(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, a.EnvConf.GetBaseURL(), http.StatusFound)
		return
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func (a *API) ProviderLogout(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, ProviderKey)
	r = r.WithContext(context.WithValue(r.Context(), providerContextKey, provider))

	session, _ := a.Store.Get(r, SessionID)
	session.Values[AuthenticatedKey] = false
	session.Values[UserIDKey] = nil
	err := session.Save(r, w)
	if err != nil {
		a.Logger.Err(err).Msg("Could not save session")
		http.Redirect(w, r, a.EnvConf.GetBaseURL(), http.StatusInternalServerError)
		return
	}

	err = gothic.Logout(w, r)
	if err != nil {
		a.Logger.Err(err).Msg("Could not logout user")
		return
	}

	a.Logger.Info().Msg("User successfully logged out")
	http.Redirect(w, r, a.EnvConf.GetBaseURL(), http.StatusFound)
}

func (a *API) ProviderCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, ProviderKey)
	r = r.WithContext(context.WithValue(r.Context(), providerContextKey, provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		a.Logger.Err(err).Msg("Error completing user authentication with gothic.")
		return
	}
	req := a.UserService.CreateUserFromGothUser(user)

	// Hash the user ID for security
	hasher := sha256.New()
	hasher.Write([]byte(user.IDToken))
	req.IDToken = []byte(hex.EncodeToString(hasher.Sum(nil)))

	// userWithOAuth, err := a.UserService.CreateUserWithOAuth(r.Context(), req)
	// if err != nil {
	// 	a.Logger.Err(err).Msg("Failed to create user with OAuth")
	// 	http.Error(w, "Failed to create user", http.StatusInternalServerError)
	// 	return
	// }
	//
	// a.Logger.Info().
	// 	Interface("user", userWithOAuth.User).
	// 	Interface("oauth", userWithOAuth.OAuthIdentity).
	// 	Msg("User created/updated successfully")
	//
	// session, _ := a.Store.Get(r, SessionID)
	// session.Values[AuthenticatedKey] = true
	// session.Values[UserIDKey] = userWithOAuth.User.ID.String()
	// err = session.Save(r, w)
	// if err != nil {
	// 	a.Logger.Err(err).Msg("Could not save session")
	// 	http.Error(w, "Session error", http.StatusInternalServerError)
	// 	return
	// }

	http.Redirect(w, r, a.EnvConf.GetBaseURL(), http.StatusFound)
}
