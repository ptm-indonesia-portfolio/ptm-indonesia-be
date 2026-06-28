package controllers

import (
	"errors"
	"strconv"
	"strings"

	"ptm-indonesia/config"
	"ptm-indonesia/helper"
	"ptm-indonesia/model"
	"ptm-indonesia/services"
	servicesContract "ptm-indonesia/services/contract"

	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	cfg         *config.AppConfig
	authService servicesContract.AuthService
	userService servicesContract.UserService
	responder   *helper.Responder
	localizer   *helper.Localizer
}

func NewUserController(
	cfg *config.AppConfig,
	authService servicesContract.AuthService,
	userService servicesContract.UserService,
	responder *helper.Responder,
	localizer *helper.Localizer,
) *UserController {
	return &UserController{
		cfg:         cfg,
		authService: authService,
		userService: userService,
		responder:   responder,
		localizer:   localizer,
	}
}

func (u *UserController) List(c fiber.Ctx) error {
	locale := u.localizer.Resolve(c.Get("Accept-Language"))

	var request model.UserListRequest
	if err := c.Bind().Query(&request); err != nil {
		return err
	}

	response, err := u.userService.List(c.Context(), request)
	if err != nil {
		if errors.Is(err, services.ErrUserInvalidSort) {
			return u.responder.Errors(c, fiber.StatusBadRequest, []string{u.localizer.MustLocalize(locale, "user.invalid_sort")})
		}

		return err
	}

	return u.responder.Success(c, fiber.StatusOK, u.localizer.MustLocalize(locale, "user.listed"), response)
}

func (u *UserController) Create(c fiber.Ctx) error {
	locale := u.localizer.Resolve(c.Get("Accept-Language"))

	currentUser, ok := helper.GetAuthUser(c)
	if !ok {
		return u.responder.Errors(c, fiber.StatusUnauthorized, []string{u.localizer.MustLocalize(locale, "error.unauthorized")})
	}

	var request model.UserCreateRequest
	if err := c.Bind().Body(&request); err != nil {
		return err
	}

	response, err := u.userService.Create(c.Context(), request, currentUser.ID)
	if err != nil {
		if errors.Is(err, services.ErrUserInvalidStatus) {
			return u.responder.Errors(c, fiber.StatusBadRequest, []string{u.localizer.MustLocalize(locale, "user.invalid_status")})
		}

		if errors.Is(err, services.ErrUserEmailAlreadyUsed) {
			return u.responder.Errors(c, fiber.StatusBadRequest, []string{u.localizer.MustLocalize(locale, "user.email_exists")})
		}

		return err
	}

	return u.responder.Success(c, fiber.StatusCreated, u.localizer.MustLocalize(locale, "user.created"), response)
}

func (u *UserController) Detail(c fiber.Ctx) error {
	locale := u.localizer.Resolve(c.Get("Accept-Language"))

	userID, err := u.parseUserID(c)
	if err != nil {
		return u.responder.Errors(c, fiber.StatusBadRequest, []string{u.localizer.MustLocalize(locale, "user.invalid_id")})
	}

	response, err := u.userService.FindByID(c.Context(), userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return u.responder.Errors(c, fiber.StatusNotFound, []string{u.localizer.MustLocalize(locale, "user.not_found")})
		}

		return err
	}

	return u.responder.Success(c, fiber.StatusOK, u.localizer.MustLocalize(locale, "user.detail"), response)
}

func (u *UserController) Update(c fiber.Ctx) error {
	locale := u.localizer.Resolve(c.Get("Accept-Language"))

	currentUser, ok := helper.GetAuthUser(c)
	if !ok {
		return u.responder.Errors(c, fiber.StatusUnauthorized, []string{u.localizer.MustLocalize(locale, "error.unauthorized")})
	}

	userID, err := u.parseUserID(c)
	if err != nil {
		return u.responder.Errors(c, fiber.StatusBadRequest, []string{u.localizer.MustLocalize(locale, "user.invalid_id")})
	}

	var request model.UserUpdateRequest
	if err := c.Bind().Body(&request); err != nil {
		return err
	}

	response, err := u.userService.Update(c.Context(), userID, request, currentUser.ID)
	if err != nil {
		if errors.Is(err, services.ErrUserInvalidStatus) {
			return u.responder.Errors(c, fiber.StatusBadRequest, []string{u.localizer.MustLocalize(locale, "user.invalid_status")})
		}

		if errors.Is(err, services.ErrUserNotFound) {
			return u.responder.Errors(c, fiber.StatusNotFound, []string{u.localizer.MustLocalize(locale, "user.not_found")})
		}

		if errors.Is(err, services.ErrUserEmailAlreadyUsed) {
			return u.responder.Errors(c, fiber.StatusBadRequest, []string{u.localizer.MustLocalize(locale, "user.email_exists")})
		}

		return err
	}

	return u.responder.Success(c, fiber.StatusOK, u.localizer.MustLocalize(locale, "user.updated"), response)
}

func (u *UserController) Delete(c fiber.Ctx) error {
	locale := u.localizer.Resolve(c.Get("Accept-Language"))

	currentUser, ok := helper.GetAuthUser(c)
	if !ok {
		return u.responder.Errors(c, fiber.StatusUnauthorized, []string{u.localizer.MustLocalize(locale, "error.unauthorized")})
	}

	userID, err := u.parseUserID(c)
	if err != nil {
		return u.responder.Errors(c, fiber.StatusBadRequest, []string{u.localizer.MustLocalize(locale, "user.invalid_id")})
	}

	if err := u.userService.Delete(c.Context(), userID, currentUser.ID); err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return u.responder.Errors(c, fiber.StatusNotFound, []string{u.localizer.MustLocalize(locale, "user.not_found")})
		}

		return err
	}

	return u.responder.Success(c, fiber.StatusOK, u.localizer.MustLocalize(locale, "user.deleted"), nil)
}

func (u *UserController) RequireLogin(c fiber.Ctx) error {
	locale := u.localizer.Resolve(c.Get("Accept-Language"))
	user, err := u.authService.GetCurrentUser(c.Context(), c.Cookies(u.cfg.Auth.CookieName))
	if err != nil {
		if errors.Is(err, services.ErrAuthMissingToken) || errors.Is(err, services.ErrAuthInvalidToken) {
			return u.responder.Errors(c, fiber.StatusUnauthorized, []string{u.localizer.MustLocalize(locale, "error.unauthorized")})
		}

		if errors.Is(err, services.ErrAuthForbidden) {
			return u.responder.Errors(c, fiber.StatusForbidden, []string{u.localizer.MustLocalize(locale, "error.forbidden_login")})
		}

		return err
	}

	if user.Status != model.UserStatusSuperAdmin {
		return u.responder.Errors(c, fiber.StatusForbidden, []string{u.localizer.MustLocalize(locale, "error.forbidden_access")})
	}

	helper.SetAuthUser(c, user)

	return c.Next()
}

func (u *UserController) parseUserID(c fiber.Ctx) (uint64, error) {
	rawUserID := strings.TrimSpace(c.Params("id"))
	if rawUserID == "" {
		return 0, errors.New("empty user id")
	}

	userID, err := strconv.ParseUint(rawUserID, 10, 64)
	if err != nil || userID == 0 {
		return 0, errors.New("invalid user id")
	}

	return userID, nil
}
