package contract

import (
	"context"

	"ptm-indonesia/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	UpdateGoogleIdentity(ctx context.Context, userID uint64, googleID string, avatarURL *string) error
}
