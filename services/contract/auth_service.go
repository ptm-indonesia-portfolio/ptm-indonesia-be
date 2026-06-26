package contract

import (
	"context"

	"ptm-indonesia/model"
)

type AuthService interface {
	PrepareGoogleLogin(ctx context.Context) (string, string, error)
	AuthenticateWithGoogle(ctx context.Context, code string, state string, expectedState string) (*model.AuthLoginResult, error)
	RefreshSession(ctx context.Context, refreshToken string) (*model.AuthLoginResult, error)
	GetCurrentUser(ctx context.Context, token string) (*model.AuthSessionUser, error)
	Logout(ctx context.Context, refreshToken string) error
}
