package contract

import (
	"context"

	"ptm-indonesia/model"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken *model.AuthRefreshToken) error
	FindActiveByTokenHash(ctx context.Context, tokenHash string) (*model.AuthRefreshToken, error)
	Rotate(ctx context.Context, currentTokenID uint64, replacement *model.AuthRefreshToken) error
	RevokeByTokenHash(ctx context.Context, tokenHash string) error
}
