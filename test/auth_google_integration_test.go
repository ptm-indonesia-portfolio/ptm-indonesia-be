//go:build integration

package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"ptm-indonesia/bootstrap"
	"ptm-indonesia/model"

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

func TestGoogleSSORejectsInactiveRegisteredUserIntegration(t *testing.T) {
	cfg, dbCleanup, dbTeardown := setupIntegrationEnvironment(t)
	defer dbCleanup()
	defer dbTeardown()

	seedIntegrationUser(t, cfg, &model.User{
		Name:      "Inactive Member",
		Email:     "inactive.member@example.com",
		Status:    model.UserStatusNotActive,
		CreatedBy: 0,
		UpdatedBy: 0,
	})

	issuerURL, closeOIDCServer := startMockGoogleOIDCServer(t, "integration-google-client", model.GoogleIdentity{
		Subject:       "google-user-3",
		Email:         "inactive.member@example.com",
		EmailVerified: true,
		Name:          "Inactive Member",
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
