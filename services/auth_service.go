package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"

	"ptm-indonesia/model"
	repositoryContract "ptm-indonesia/repository/contract"
	servicesContract "ptm-indonesia/services/contract"

	"gorm.io/gorm"
)

var (
	ErrAuthInvalidState        = errors.New("invalid auth state")
	ErrAuthMissingCode         = errors.New("missing auth code")
	ErrAuthForbidden           = errors.New("forbidden login")
	ErrAuthMissingToken        = errors.New("missing session token")
	ErrAuthInvalidToken        = errors.New("invalid session token")
	ErrAuthMissingRefreshToken = errors.New("missing refresh token")
	ErrAuthInvalidRefreshToken = errors.New("invalid refresh token")
)

type AuthService struct {
	googleOIDCClient       servicesContract.GoogleOIDCClient
	sessionTokenManager    servicesContract.SessionTokenManager
	refreshTokenManager    servicesContract.RefreshTokenManager
	userRepository         repositoryContract.UserRepository
	refreshTokenRepository repositoryContract.RefreshTokenRepository
}

func NewAuthService(
	googleOIDCClient servicesContract.GoogleOIDCClient,
	sessionTokenManager servicesContract.SessionTokenManager,
	refreshTokenManager servicesContract.RefreshTokenManager,
	userRepository repositoryContract.UserRepository,
	refreshTokenRepository repositoryContract.RefreshTokenRepository,
) *AuthService {
	return &AuthService{
		googleOIDCClient:       googleOIDCClient,
		sessionTokenManager:    sessionTokenManager,
		refreshTokenManager:    refreshTokenManager,
		userRepository:         userRepository,
		refreshTokenRepository: refreshTokenRepository,
	}
}

func (s *AuthService) PrepareGoogleLogin(_ context.Context) (string, string, error) {
	state, err := generateAuthState()
	if err != nil {
		return "", "", err
	}

	return state, s.googleOIDCClient.AuthCodeURL(state), nil
}

func (s *AuthService) AuthenticateWithGoogle(ctx context.Context, code string, state string, expectedState string) (*model.AuthLoginResult, error) {
	if strings.TrimSpace(code) == "" {
		return nil, ErrAuthMissingCode
	}

	if strings.TrimSpace(state) == "" || strings.TrimSpace(expectedState) == "" || state != expectedState {
		return nil, ErrAuthInvalidState
	}

	identity, err := s.googleOIDCClient.ExchangeCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if !identity.EmailVerified {
		return nil, ErrAuthForbidden
	}

	user, err := s.userRepository.FindByEmail(ctx, identity.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuthForbidden
		}

		return nil, err
	}

	if !model.IsLoginAllowedUserStatus(user.Status) {
		return nil, ErrAuthForbidden
	}

	if err := s.userRepository.UpdateGoogleIdentity(ctx, user.ID, identity.Subject, identity.Picture); err != nil {
		return nil, err
	}

	user.GoogleID = &identity.Subject
	user.AvatarURL = identity.Picture

	accessToken, err := s.sessionTokenManager.Generate(user)
	if err != nil {
		return nil, err
	}

	refreshTokenModel, refreshToken, err := s.refreshTokenManager.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepository.Create(ctx, refreshTokenModel); err != nil {
		return nil, err
	}

	return &model.AuthLoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         toAuthSessionUser(user),
	}, nil
}

func (s *AuthService) RefreshSession(ctx context.Context, refreshToken string) (*model.AuthLoginResult, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, ErrAuthMissingRefreshToken
	}

	currentRefreshToken, err := s.refreshTokenRepository.FindActiveByTokenHash(ctx, s.refreshTokenManager.Hash(refreshToken))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuthInvalidRefreshToken
		}

		return nil, err
	}

	user, err := s.userRepository.FindByID(ctx, currentRefreshToken.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuthInvalidRefreshToken
		}

		return nil, err
	}

	if !model.IsLoginAllowedUserStatus(user.Status) {
		return nil, ErrAuthForbidden
	}

	accessToken, err := s.sessionTokenManager.Generate(user)
	if err != nil {
		return nil, err
	}

	replacementRefreshToken, plainRefreshToken, err := s.refreshTokenManager.Generate(user.ID)
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepository.Rotate(ctx, currentRefreshToken.ID, replacementRefreshToken); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuthInvalidRefreshToken
		}

		return nil, err
	}

	return &model.AuthLoginResult{
		AccessToken:  accessToken,
		RefreshToken: plainRefreshToken,
		User:         toAuthSessionUser(user),
	}, nil
}

func (s *AuthService) GetCurrentUser(ctx context.Context, token string) (*model.AuthSessionUser, error) {
	if strings.TrimSpace(token) == "" {
		return nil, ErrAuthMissingToken
	}

	sessionUser, err := s.sessionTokenManager.Parse(token)
	if err != nil {
		return nil, ErrAuthInvalidToken
	}

	user, err := s.userRepository.FindByID(ctx, sessionUser.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuthInvalidToken
		}

		return nil, err
	}

	if !model.IsLoginAllowedUserStatus(user.Status) {
		return nil, ErrAuthForbidden
	}

	return toAuthSessionUser(user), nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return nil
	}

	return s.refreshTokenRepository.RevokeByTokenHash(ctx, s.refreshTokenManager.Hash(refreshToken))
}

func generateAuthState() (string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func toAuthSessionUser(user *model.User) *model.AuthSessionUser {
	return &model.AuthSessionUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		GoogleID:  user.GoogleID,
		AvatarURL: user.AvatarURL,
		Address:   user.Address,
		Telp:      user.Telp,
		Status:    user.Status,
	}
}
