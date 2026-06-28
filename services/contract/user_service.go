package contract

import (
	"context"

	"ptm-indonesia/model"
)

type UserService interface {
	List(ctx context.Context, request model.UserListRequest) (*model.UserListResponse, error)
	FindByID(ctx context.Context, id uint64) (*model.UserResponse, error)
	Create(ctx context.Context, request model.UserCreateRequest, actorID uint64) (*model.UserResponse, error)
	Update(ctx context.Context, id uint64, request model.UserUpdateRequest, actorID uint64) (*model.UserResponse, error)
	Delete(ctx context.Context, id uint64, actorID uint64) error
}
