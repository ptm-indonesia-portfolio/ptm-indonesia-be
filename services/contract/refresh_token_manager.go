package contract

import "ptm-indonesia/model"

type RefreshTokenManager interface {
	Generate(userID uint64) (*model.AuthRefreshToken, string, error)
	Hash(token string) string
}
