package services

import (
	"context"
	"errors"
	"testing"

	"ptm-indonesia/config"

	"github.com/stretchr/testify/assert"
)

type stubSystemRepository struct {
	err error
}

func (s *stubSystemRepository) Ping(_ context.Context) error {
	return s.err
}

func TestHealthServiceCheckSuccess(t *testing.T) {
	t.Parallel()

	service := NewHealthService(&config.AppConfig{
		App: config.AppSection{
			Name:              "PTM Indonesia API",
			Environment:       "development",
			Timezone:          "Asia/Jakarta",
			DefaultLanguage:   "id",
			SupportedLanguage: []string{"id", "en"},
		},
	}, &stubSystemRepository{})

	response := service.Check(context.Background())

	assert.Equal(t, "PTM Indonesia API", response.Name)
	assert.Equal(t, "development", response.Environment)
	assert.Equal(t, "up", response.Database)
	assert.Equal(t, "id", response.DefaultLanguage)
	assert.Equal(t, []string{"id", "en"}, response.SupportedLanguages)
	assert.NotEmpty(t, response.Timestamp)
}

func TestHealthServiceCheckDatabaseDown(t *testing.T) {
	t.Parallel()

	service := NewHealthService(&config.AppConfig{
		App: config.AppSection{
			Name:              "PTM Indonesia API",
			Environment:       "production",
			Timezone:          "Asia/Jakarta",
			DefaultLanguage:   "en",
			SupportedLanguage: []string{"id", "en"},
		},
	}, &stubSystemRepository{
		err: errors.New("database unavailable"),
	})

	response := service.Check(context.Background())

	assert.Equal(t, "down", response.Database)
	assert.Equal(t, "production", response.Environment)
}
