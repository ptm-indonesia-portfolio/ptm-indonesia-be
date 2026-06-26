package contract

import "ptm-indonesia/model"

type SessionTokenManager interface {
	Generate(user *model.User) (string, error)
	Parse(token string) (*model.AuthSessionUser, error)
}
