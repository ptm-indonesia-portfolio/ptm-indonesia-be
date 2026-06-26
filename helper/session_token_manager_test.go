package helper

import (
	"testing"
	"time"

	"ptm-indonesia/config"
	"ptm-indonesia/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionTokenManagerGenerateAndParse(t *testing.T) {
	t.Parallel()

	tokenManager := NewSessionTokenManager(&config.AppConfig{
		App: config.AppSection{
			Name: "PTM Indonesia API",
		},
		Auth: config.AuthSection{
			CookieSecret:   "test-secret",
			AccessTokenTTL: 1 * time.Hour,
		},
	})

	googleID := "google-123"
	token, err := tokenManager.Generate(&model.User{
		ID:       10,
		Name:     "Super Admin",
		Email:    "admin@example.com",
		GoogleID: &googleID,
		Status:   model.UserStatusSuperAdmin,
	})
	require.NoError(t, err)

	sessionUser, err := tokenManager.Parse(token)
	require.NoError(t, err)

	assert.Equal(t, uint64(10), sessionUser.ID)
	assert.Equal(t, "Super Admin", sessionUser.Name)
	assert.Equal(t, "admin@example.com", sessionUser.Email)
	assert.Equal(t, model.UserStatusSuperAdmin, sessionUser.Status)
	require.NotNil(t, sessionUser.GoogleID)
	assert.Equal(t, "google-123", *sessionUser.GoogleID)
}
