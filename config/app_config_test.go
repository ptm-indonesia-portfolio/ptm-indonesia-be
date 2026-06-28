package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAppConfigUsesEnvironmentVariablesWithoutDotEnvFile(t *testing.T) {
	t.Setenv("APP_PORT", "9100")
	t.Setenv("DB_HOST", "ptm_indonesia_be_db")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "ptm_indonesia")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "postgres")
	t.Setenv("EMAIL_ADMIN", "admin@example.com")

	cfg, err := NewAppConfig()
	require.NoError(t, err)

	assert.Equal(t, "9100", cfg.App.Port)
	assert.Equal(t, "ptm_indonesia_be_db", cfg.Database.Host)
}

func TestNewAppConfigReturnsErrorWhenExplicitConfigFileIsMissing(t *testing.T) {
	t.Setenv("EMAIL_ADMIN", "admin@example.com")
	t.Setenv("ENV_FILE_PATH", filepath.Join(t.TempDir(), "missing.env"))

	_, err := NewAppConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "config file not found")
}

func TestNewAppConfigIncludesFrontendOriginInCORSOrigins(t *testing.T) {
	t.Setenv("EMAIL_ADMIN", "admin@example.com")
	t.Setenv("CORS_ORIGINS", "http://localhost:3100")
	t.Setenv("FRONTEND_URL", "http://localhost:3101/home")

	cfg, err := NewAppConfig()
	require.NoError(t, err)

	assert.Contains(t, cfg.App.CORSOrigins, "http://localhost:3100")
	assert.Contains(t, cfg.App.CORSOrigins, "http://localhost:3101")
}

func TestNewAppConfigReturnsErrorWhenEmailAdminIsMissing(t *testing.T) {
	t.Setenv("EMAIL_ADMIN", "   ")

	_, err := NewAppConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "EMAIL_ADMIN is required")
}
