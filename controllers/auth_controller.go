package controllers

import (
	"errors"
	"net/url"
	"time"

	"ptm-indonesia/config"
	"ptm-indonesia/helper"
	"ptm-indonesia/model"
	"ptm-indonesia/services"
	servicesContract "ptm-indonesia/services/contract"

	"github.com/gofiber/fiber/v3"
)

const googleAuthStateCookieName = "ptm_google_auth_state"

type AuthController struct {
	cfg         *config.AppConfig
	authService servicesContract.AuthService
	responder   *helper.Responder
	localizer   *helper.Localizer
}

func NewAuthController(
	cfg *config.AppConfig,
	authService servicesContract.AuthService,
	responder *helper.Responder,
	localizer *helper.Localizer,
) *AuthController {
	return &AuthController{
		cfg:         cfg,
		authService: authService,
		responder:   responder,
		localizer:   localizer,
	}
}

func (a *AuthController) GoogleLogin(c fiber.Ctx) error {
	locale := a.localizer.Resolve(c.Get("Accept-Language"))
	state, loginURL, err := a.authService.PrepareGoogleLogin(c.Context())
	if err != nil {
		return err
	}

	a.setStateCookie(c, state)

	return a.responder.Success(c, fiber.StatusOK, a.localizer.MustLocalize(locale, "auth.google_login_url"), model.AuthGoogleLoginResponse{
		URL: loginURL,
	})
}

func (a *AuthController) GoogleCallback(c fiber.Ctx) error {
	expectedState := c.Cookies(googleAuthStateCookieName)
	loginResult, err := a.authService.AuthenticateWithGoogle(
		c.Context(),
		c.Query("code"),
		c.Query("state"),
		expectedState,
	)
	if err != nil {
		a.clearStateCookie(c)
		return a.handleGoogleCallbackError(c, err)
	}

	a.clearStateCookie(c)
	a.setAuthCookie(c, loginResult.AccessToken)
	a.setRefreshCookie(c, loginResult.RefreshToken)
	a.setLoggedInCookie(c)

	redirectURL := buildFrontendRedirectURL(a.cfg.Auth.FrontendURL, "success")
	return c.Redirect().To(redirectURL)
}

func (a *AuthController) Refresh(c fiber.Ctx) error {
	locale := a.localizer.Resolve(c.Get("Accept-Language"))
	refreshResult, err := a.authService.RefreshSession(c.Context(), c.Cookies(a.cfg.Auth.RefreshCookieName))
	if err != nil {
		if errors.Is(err, services.ErrAuthMissingRefreshToken) || errors.Is(err, services.ErrAuthInvalidRefreshToken) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{
				Errors: []string{a.localizer.MustLocalize(locale, "error.unauthorized")},
			})
		}

		if errors.Is(err, services.ErrAuthForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
				Errors: []string{a.localizer.MustLocalize(locale, "error.forbidden_login")},
			})
		}

		return err
	}

	a.setAuthCookie(c, refreshResult.AccessToken)
	a.setRefreshCookie(c, refreshResult.RefreshToken)
	a.setLoggedInCookie(c)

	return a.responder.Success(c, fiber.StatusOK, a.localizer.MustLocalize(locale, "auth.refreshed"), refreshResult.User)
}

func (a *AuthController) Me(c fiber.Ctx) error {
	locale := a.localizer.Resolve(c.Get("Accept-Language"))
	user, err := a.authService.GetCurrentUser(c.Context(), c.Cookies(a.cfg.Auth.CookieName))
	if err != nil {
		if errors.Is(err, services.ErrAuthMissingToken) || errors.Is(err, services.ErrAuthInvalidToken) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{
				Errors: []string{a.localizer.MustLocalize(locale, "error.unauthorized")},
			})
		}

		if errors.Is(err, services.ErrAuthForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(model.ErrorResponse{
				Errors: []string{a.localizer.MustLocalize(locale, "error.forbidden_login")},
			})
		}

		return err
	}

	return a.responder.Success(c, fiber.StatusOK, a.localizer.MustLocalize(locale, "auth.current_user"), user)
}

func (a *AuthController) Logout(c fiber.Ctx) error {
	locale := a.localizer.Resolve(c.Get("Accept-Language"))
	if err := a.authService.Logout(c.Context(), c.Cookies(a.cfg.Auth.RefreshCookieName)); err != nil {
		return err
	}

	a.clearAuthCookie(c)
	a.clearRefreshCookie(c)
	a.clearLoggedInCookie(c)

	return a.responder.Success(c, fiber.StatusOK, a.localizer.MustLocalize(locale, "auth.logged_out"), nil)
}

func (a *AuthController) setStateCookie(c fiber.Ctx, state string) {
	c.Cookie(&fiber.Cookie{
		Name:     googleAuthStateCookieName,
		Value:    state,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Secure:   true,
		MaxAge:   600,
		Path:     "/",
	})
}

func (a *AuthController) clearStateCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     googleAuthStateCookieName,
		Value:    "",
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
	})
}

func (a *AuthController) setAuthCookie(c fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     a.cfg.Auth.CookieName,
		Value:    token,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Secure:   true,
		MaxAge:   int(a.cfg.Auth.AccessTokenTTL.Seconds()),
		Path:     "/",
	})
}

func (a *AuthController) setRefreshCookie(c fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     a.cfg.Auth.RefreshCookieName,
		Value:    token,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Secure:   true,
		MaxAge:   int(a.cfg.Auth.RefreshTokenTTL.Seconds()),
		Path:     "/",
	})
}

func (a *AuthController) setLoggedInCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     a.cfg.Auth.LoggedInCookieName,
		Value:    "1",
		HTTPOnly: false,
		SameSite: fiber.CookieSameSiteLaxMode,
		Secure:   false,
		MaxAge:   int(a.cfg.Auth.RefreshTokenTTL.Seconds()),
		Path:     "/",
	})
}

func (a *AuthController) clearAuthCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     a.cfg.Auth.CookieName,
		Value:    "",
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
	})
}

func (a *AuthController) clearRefreshCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     a.cfg.Auth.RefreshCookieName,
		Value:    "",
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
	})
}

func (a *AuthController) clearLoggedInCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     a.cfg.Auth.LoggedInCookieName,
		Value:    "",
		HTTPOnly: false,
		SameSite: fiber.CookieSameSiteLaxMode,
		Secure:   false,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Path:     "/",
	})
}

func (a *AuthController) handleGoogleCallbackError(c fiber.Ctx, err error) error {
	status := "error"
	if errors.Is(err, services.ErrAuthForbidden) {
		status = "forbidden"
	}

	if errors.Is(err, services.ErrAuthInvalidState) || errors.Is(err, services.ErrAuthMissingCode) {
		status = "invalid_request"
	}

	return c.Redirect().To(buildFrontendRedirectURL(a.cfg.Auth.FrontendURL, status))
}

func buildFrontendRedirectURL(frontendURL string, status string) string {
	if frontendURL == "" {
		return "/"
	}

	parsedURL, err := url.Parse(frontendURL)
	if err != nil {
		return frontendURL
	}

	query := parsedURL.Query()
	query.Set("auth_status", status)
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String()
}
