package helper

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"ptm-indonesia/config"
	"ptm-indonesia/model"
)

type RefreshTokenManager struct {
	refreshTokenTTL time.Duration
}

func NewRefreshTokenManager(cfg *config.AppConfig) *RefreshTokenManager {
	return &RefreshTokenManager{
		refreshTokenTTL: cfg.Auth.RefreshTokenTTL,
	}
}

func (m *RefreshTokenManager) Generate(userID uint64) (*model.AuthRefreshToken, string, error) {
	buffer := make([]byte, 32)
	if _, err := rand.Read(buffer); err != nil {
		return nil, "", err
	}

	token := base64.RawURLEncoding.EncodeToString(buffer)

	return &model.AuthRefreshToken{
		UserID:    userID,
		TokenHash: m.Hash(token),
		ExpiresAt: time.Now().Add(m.refreshTokenTTL),
		CreatedBy: 0,
		UpdatedBy: 0,
	}, token, nil
}

func (m *RefreshTokenManager) Hash(token string) string {
	hash := sha256.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}
