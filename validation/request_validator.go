package validation

import (
	"fmt"
	"reflect"
	"strings"

	"ptm-indonesia/config"

	enlocale "github.com/go-playground/locales/en"
	idlocale "github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	idtranslations "github.com/go-playground/validator/v10/translations/id"
)

type RequestValidator struct {
	validate      *validator.Validate
	translators   map[string]ut.Translator
	defaultLocale string
}

func NewRequestValidator(cfg *config.AppConfig) (*RequestValidator, error) {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "" || name == "-" {
			return strings.ToLower(field.Name)
		}

		return name
	})

	english := enlocale.New()
	indonesian := idlocale.New()
	uni := ut.New(english, english, indonesian)

	enTranslator, found := uni.GetTranslator("en")
	if !found {
		return nil, fmt.Errorf("english translator not found")
	}

	idTranslator, found := uni.GetTranslator("id")
	if !found {
		return nil, fmt.Errorf("indonesian translator not found")
	}

	if err := entranslations.RegisterDefaultTranslations(validate, enTranslator); err != nil {
		return nil, fmt.Errorf("register english translations: %w", err)
	}

	if err := idtranslations.RegisterDefaultTranslations(validate, idTranslator); err != nil {
		return nil, fmt.Errorf("register indonesian translations: %w", err)
	}

	if err := registerRequiredTranslation(validate, enTranslator, "The {0} is required!"); err != nil {
		return nil, fmt.Errorf("register english required translation: %w", err)
	}

	if err := registerRequiredTranslation(validate, idTranslator, "{0} wajib diisi!"); err != nil {
		return nil, fmt.Errorf("register indonesian required translation: %w", err)
	}

	return &RequestValidator{
		validate:      validate,
		translators:   map[string]ut.Translator{"en": enTranslator, "id": idTranslator},
		defaultLocale: cfg.App.DefaultLanguage,
	}, nil
}

func (v *RequestValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func (v *RequestValidator) TranslateValidationErrors(locale string, validationErrors validator.ValidationErrors) []string {
	translator, ok := v.translators[strings.ToLower(locale)]
	if !ok {
		translator = v.translators[v.defaultLocale]
	}

	errors := make([]string, 0, len(validationErrors))
	for _, fieldError := range validationErrors {
		errors = append(errors, fieldError.Translate(translator))
	}

	return errors
}

func registerRequiredTranslation(validate *validator.Validate, translator ut.Translator, message string) error {
	return validate.RegisterTranslation(
		"required",
		translator,
		func(translation ut.Translator) error {
			return translation.Add("required", message, true)
		},
		func(translation ut.Translator, fieldError validator.FieldError) string {
			result, err := translation.T("required", fieldError.Field())
			if err != nil {
				return fieldError.Error()
			}

			return result
		},
	)
}
