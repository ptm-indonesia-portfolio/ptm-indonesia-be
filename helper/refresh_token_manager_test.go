package helper

import (
	"testing"
	"time"

	"ptm-indonesia/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshTokenManagerGenerateAndHash(t *testing.T) {
	t.Parallel()

	tokenManager := NewRefreshTokenManager(&config.AppConfig{
		Auth: config.AuthSection{
			RefreshTokenTTL: 24 * time.Hour,
		},
	})

	refreshToken, plainToken, err := tokenManager.Generate(42)
	require.NoError(t, err)

	require.NotEmpty(t, plainToken)
	assert.Equal(t, uint64(42), refreshToken.UserID)
	assert.Equal(t, tokenManager.Hash(plainToken), refreshToken.TokenHash)
	assert.True(t, refreshToken.ExpiresAt.After(time.Now()))
}
