package validation

import (
	"testing"

	"ptm-indonesia/config"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type createUserRequest struct {
	Email string `json:"email" validate:"required"`
}

func TestRequestValidatorTranslateValidationErrorsEnglish(t *testing.T) {
	t.Parallel()

	requestValidator, err := NewRequestValidator(&config.AppConfig{
		App: config.AppSection{
			DefaultLanguage: "en",
		},
	})
	require.NoError(t, err)

	err = requestValidator.Validate(createUserRequest{})
	require.Error(t, err)

	validationErrors, ok := err.(validator.ValidationErrors)
	require.True(t, ok)

	errors := requestValidator.TranslateValidationErrors("en", validationErrors)
	assert.Equal(t, []string{"The email is required!"}, errors)
}

func TestRequestValidatorTranslateValidationErrorsIndonesian(t *testing.T) {
	t.Parallel()

	requestValidator, err := NewRequestValidator(&config.AppConfig{
		App: config.AppSection{
			DefaultLanguage: "id",
		},
	})
	require.NoError(t, err)

	err = requestValidator.Validate(createUserRequest{})
	require.Error(t, err)

	validationErrors, ok := err.(validator.ValidationErrors)
	require.True(t, ok)

	errors := requestValidator.TranslateValidationErrors("id", validationErrors)
	assert.Equal(t, []string{"email wajib diisi!"}, errors)
}
