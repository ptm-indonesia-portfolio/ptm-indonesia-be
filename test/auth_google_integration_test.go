//go:build integration

package test

import (
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"ptm-indonesia/bootstrap"
	"ptm-indonesia/config"
	"ptm-indonesia/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestGoogleSSOLoginAndRefreshIntegration(t *testing.T) {
	cfg, dbCleanup, dbTeardown := setupIntegrationEnvironment(t)
	defer dbCleanup()
	defer dbTeardown()

	issuerURL, closeOIDCServer := startMockGoogleOIDCServer(t, "integration-google-client", model.GoogleIdentity{
		Subject:       "google-user-1",
		Email:         cfg.Admin.Email,
		EmailVerified: true,
		Name:          "Google Admin",
		Picture:       stringPointer("https://example.com/avatar.png"),
	})
	defer closeOIDCServer()

	t.Setenv("GOOGLE_CLIENT_ID", "integration-google-client")
	t.Setenv("GOOGLE_CLIENT_SECRET", "integration-google-secret")
	t.Setenv("GOOGLE_ISSUER_URL", issuerURL)
	t.Setenv("GOOGLE_REDIRECT_URI", "http://localhost:9101/api/v1/auth/google/callback")
	t.Setenv("FRONTEND_URL", "http://localhost:3100/home")

	app, cleanup, err := bootstrap.InitializeHTTPApplication()
	require.NoError(t, err)
	defer cleanup()

	stateCookie := prepareGoogleLoginState(t, app)

	callbackRequest := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/auth/google/callback?code=test-auth-code&state="+url.QueryEscape(stateCookie.Value),
		nil,
	)
	callbackRequest.AddCookie(stateCookie)

	callbackResponse, err := app.App.Test(callbackRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusSeeOther, callbackResponse.StatusCode)
	require.Contains(t, callbackResponse.Header.Get("Location"), "auth_status=success")

	authCookie := findCookieByName(callbackResponse.Cookies(), cfg.Auth.CookieName)
	refreshCookie := findCookieByName(callbackResponse.Cookies(), cfg.Auth.RefreshCookieName)
	loggedInCookie := findCookieByName(callbackResponse.Cookies(), cfg.Auth.LoggedInCookieName)

	require.NotNil(t, authCookie)
	require.True(t, authCookie.HttpOnly)
	require.True(t, authCookie.Secure)
	require.Equal(t, int(cfg.Auth.AccessTokenTTL.Seconds()), authCookie.MaxAge)

	require.NotNil(t, refreshCookie)
	require.True(t, refreshCookie.HttpOnly)
	require.True(t, refreshCookie.Secure)
	require.Equal(t, int(cfg.Auth.RefreshTokenTTL.Seconds()), refreshCookie.MaxAge)

	require.NotNil(t, loggedInCookie)
	require.False(t, loggedInCookie.HttpOnly)
	require.False(t, loggedInCookie.Secure)
	require.Equal(t, "1", loggedInCookie.Value)
	require.Equal(t, int(cfg.Auth.RefreshTokenTTL.Seconds()), loggedInCookie.MaxAge)

	mePayload := requestCurrentUser(t, app, authCookie)
	require.Equal(t, cfg.Admin.Email, mePayload.Data.Email)
	require.Equal(t, model.UserStatusSuperAdmin, mePayload.Data.Status)
	require.NotNil(t, mePayload.Data.GoogleID)

	refreshRequest := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
	refreshRequest.AddCookie(refreshCookie)

	refreshResponse, err := app.App.Test(refreshRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, refreshResponse.StatusCode)

	var refreshPayload struct {
		Message string                `json:"message"`
		Data    model.AuthSessionUser `json:"data"`
	}
	require.NoError(t, json.NewDecoder(refreshResponse.Body).Decode(&refreshPayload))
	require.Equal(t, cfg.Admin.Email, refreshPayload.Data.Email)

	refreshedAuthCookie := findCookieByName(refreshResponse.Cookies(), cfg.Auth.CookieName)
	refreshedRefreshCookie := findCookieByName(refreshResponse.Cookies(), cfg.Auth.RefreshCookieName)
	refreshedLoggedInCookie := findCookieByName(refreshResponse.Cookies(), cfg.Auth.LoggedInCookieName)

	require.NotNil(t, refreshedAuthCookie)
	require.True(t, refreshedAuthCookie.HttpOnly)
	require.True(t, refreshedAuthCookie.Secure)
	require.Equal(t, int(cfg.Auth.AccessTokenTTL.Seconds()), refreshedAuthCookie.MaxAge)

	require.NotNil(t, refreshedRefreshCookie)
	require.True(t, refreshedRefreshCookie.HttpOnly)
	require.True(t, refreshedRefreshCookie.Secure)
	require.NotEqual(t, refreshCookie.Value, refreshedRefreshCookie.Value)

	require.NotNil(t, refreshedLoggedInCookie)
	require.False(t, refreshedLoggedInCookie.HttpOnly)
	require.False(t, refreshedLoggedInCookie.Secure)
	require.Equal(t, "1", refreshedLoggedInCookie.Value)

	mePayload = requestCurrentUser(t, app, refreshedAuthCookie)
	require.Equal(t, cfg.Admin.Email, mePayload.Data.Email)
}

func TestGoogleSSORejectsUnregisteredEmailIntegration(t *testing.T) {
	cfg, dbCleanup, dbTeardown := setupIntegrationEnvironment(t)
	defer dbCleanup()
	defer dbTeardown()

	issuerURL, closeOIDCServer := startMockGoogleOIDCServer(t, "integration-google-client", model.GoogleIdentity{
		Subject:       "google-user-2",
		Email:         "stranger@example.com",
		EmailVerified: true,
		Name:          "Unknown User",
	})
	defer closeOIDCServer()

	t.Setenv("GOOGLE_CLIENT_ID", "integration-google-client")
	t.Setenv("GOOGLE_CLIENT_SECRET", "integration-google-secret")
	t.Setenv("GOOGLE_ISSUER_URL", issuerURL)
	t.Setenv("GOOGLE_REDIRECT_URI", "http://localhost:9101/api/v1/auth/google/callback")
	t.Setenv("FRONTEND_URL", "http://localhost:3100/home")

	app, cleanup, err := bootstrap.InitializeHTTPApplication()
	require.NoError(t, err)
	defer cleanup()

	stateCookie := prepareGoogleLoginState(t, app)

	callbackRequest := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/auth/google/callback?code=test-auth-code&state="+url.QueryEscape(stateCookie.Value),
		nil,
	)
	callbackRequest.AddCookie(stateCookie)

	callbackResponse, err := app.App.Test(callbackRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusSeeOther, callbackResponse.StatusCode)
	require.Contains(t, callbackResponse.Header.Get("Location"), "auth_status=forbidden")
	require.Nil(t, findCookieByName(callbackResponse.Cookies(), cfg.Auth.CookieName))
	require.Nil(t, findCookieByName(callbackResponse.Cookies(), cfg.Auth.RefreshCookieName))
	require.Nil(t, findCookieByName(callbackResponse.Cookies(), cfg.Auth.LoggedInCookieName))
}

func setupIntegrationEnvironment(t *testing.T) (*config.AppConfig, func(), func()) {
	t.Helper()

	workingDirectory, err := os.Getwd()
	require.NoError(t, err)

	projectRoot := filepath.Clean(filepath.Join(workingDirectory, ".."))
	require.NoError(t, os.Chdir(projectRoot))
	t.Cleanup(func() {
		_ = os.Chdir(workingDirectory)
	})

	envFilePath := filepath.Join(projectRoot, ".env.testing")
	migrationPath := filepath.Join(projectRoot, "db", "migrations")

	t.Setenv("ENV_FILE_PATH", envFilePath)
	t.Setenv("MIGRATION_SOURCE", "file://"+filepath.ToSlash(migrationPath))

	return setupIntegrationDatabase(t)
}

func prepareGoogleLoginState(t *testing.T, app *bootstrap.HTTPApplication) *http.Cookie {
	t.Helper()

	loginRequest := httptest.NewRequest(http.MethodGet, "/api/v1/auth/google/login", nil)
	loginResponse, err := app.App.Test(loginRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, loginResponse.StatusCode)

	var loginPayload struct {
		Message string `json:"message"`
		Data    struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	require.NoError(t, json.NewDecoder(loginResponse.Body).Decode(&loginPayload))
	require.NotEmpty(t, loginPayload.Data.URL)

	stateCookie := findCookieByName(loginResponse.Cookies(), "ptm_google_auth_state")
	require.NotNil(t, stateCookie)
	require.True(t, stateCookie.HttpOnly)
	require.True(t, stateCookie.Secure)

	return stateCookie
}

func requestCurrentUser(t *testing.T, app *bootstrap.HTTPApplication, authCookie *http.Cookie) struct {
	Message string                `json:"message"`
	Data    model.AuthSessionUser `json:"data"`
} {
	t.Helper()

	meRequest := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	meRequest.AddCookie(authCookie)

	meResponse, err := app.App.Test(meRequest)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, meResponse.StatusCode)

	var mePayload struct {
		Message string                `json:"message"`
		Data    model.AuthSessionUser `json:"data"`
	}
	require.NoError(t, json.NewDecoder(meResponse.Body).Decode(&mePayload))

	return mePayload
}

func findCookieByName(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}

	return nil
}

func setupIntegrationDatabase(t *testing.T) (*config.AppConfig, func(), func()) {
	t.Helper()

	cfg, err := config.NewAppConfig()
	require.NoError(t, err)

	adminDB, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.SSLMode,
	))
	require.NoError(t, err)

	require.NoError(t, adminDB.Ping())
	_, err = adminDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.Database.Name))
	require.NoError(t, err)
	_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Database.Name))
	require.NoError(t, err)

	migrationConfig := config.NewMigrationRunner(cfg)
	migrator, err := migrate.New(migrationConfig.Source, migrationConfig.DatabaseURL)
	require.NoError(t, err)
	require.NoError(t, migrator.Up())
	_, _ = migrator.Close()

	gormDB, cleanup, err := config.NewDatabase(cfg)
	require.NoError(t, err)
	require.NoError(t, gormDB.Create(&model.User{
		Name:      "Super Admin",
		Email:     cfg.Admin.Email,
		Status:    model.UserStatusSuperAdmin,
		CreatedBy: 0,
		UpdatedBy: 0,
	}).Error)

	dbCleanup := func() {
		cleanup()
		_ = adminDB.Close()
	}

	dbTeardown := func() {
		teardownDB, openErr := sql.Open("postgres", fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.SSLMode,
		))
		if openErr != nil {
			return
		}
		defer teardownDB.Close()

		_, _ = teardownDB.Exec(`
			SELECT pg_terminate_backend(pid)
			FROM pg_stat_activity
			WHERE datname = $1 AND pid <> pg_backend_pid()
		`, cfg.Database.Name)
		_, _ = teardownDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.Database.Name))
	}

	return cfg, dbCleanup, dbTeardown
}

func startMockGoogleOIDCServer(t *testing.T, clientID string, identity model.GoogleIdentity) (string, func()) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	jwk := map[string]any{
		"kty": "RSA",
		"use": "sig",
		"alg": "RS256",
		"kid": "integration-key",
		"n":   base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.N.Bytes()),
		"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(privateKey.PublicKey.E)).Bytes()),
	}

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"issuer":                                server.URL,
			"authorization_endpoint":                server.URL + "/auth",
			"token_endpoint":                        server.URL + "/token",
			"jwks_uri":                              server.URL + "/jwks",
			"id_token_signing_alg_values_supported": []string{"RS256"},
		})
	})

	mux.HandleFunc("/jwks", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"keys": []any{jwk},
		})
	})

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, server.URL+"/noop", http.StatusFound)
	})

	mux.HandleFunc("/token", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		now := time.Now()
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"iss":            server.URL,
			"sub":            identity.Subject,
			"aud":            clientID,
			"exp":            now.Add(time.Hour).Unix(),
			"iat":            now.Unix(),
			"email":          identity.Email,
			"email_verified": identity.EmailVerified,
			"name":           identity.Name,
			"picture":        dereferenceString(identity.Picture),
		})
		token.Header["kid"] = "integration-key"

		signedToken, signErr := token.SignedString(privateKey)
		require.NoError(t, signErr)

		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "integration-access-token",
			"token_type":   "Bearer",
			"expires_in":   3600,
			"id_token":     signedToken,
		})
	})

	return server.URL, server.Close
}

func stringPointer(value string) *string {
	return &value
}

func dereferenceString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
