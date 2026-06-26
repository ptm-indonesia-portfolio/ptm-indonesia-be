package config

import (
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func NewI18nBundle() (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.Indonesian)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	if _, err := bundle.LoadMessageFile("lang/en.json"); err != nil {
		return nil, fmt.Errorf("load english messages: %w", err)
	}

	if _, err := bundle.LoadMessageFile("lang/id.json"); err != nil {
		return nil, fmt.Errorf("load indonesian messages: %w", err)
	}

	return bundle, nil
}
