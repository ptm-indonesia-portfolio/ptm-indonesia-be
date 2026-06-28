package contract

import (
	"context"

	"ptm-indonesia/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uint64) (*model.User, error)
	List(ctx context.Context, request model.UserListRequest, excludedEmail string) ([]model.User, error)
	Count(ctx context.Context, request model.UserListRequest, excludedEmail string) (int64, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	SoftDelete(ctx context.Context, id uint64, updatedBy uint64) error
	ExistsByEmail(ctx context.Context, email string, excludedID *uint64) (bool, error)
	UpdateGoogleIdentity(ctx context.Context, userID uint64, googleID string, avatarURL *string) error
}
