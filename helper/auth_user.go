package helper

import (
	"ptm-indonesia/model"

	"github.com/gofiber/fiber/v3"
)

const authUserContextKey = "auth_user"

func SetAuthUser(c fiber.Ctx, user *model.AuthSessionUser) {
	c.Locals(authUserContextKey, user)
}

func GetAuthUser(c fiber.Ctx) (*model.AuthSessionUser, bool) {
	value := c.Locals(authUserContextKey)
	user, ok := value.(*model.AuthSessionUser)
	if !ok || user == nil {
		return nil, false
	}

	return user, true
}
