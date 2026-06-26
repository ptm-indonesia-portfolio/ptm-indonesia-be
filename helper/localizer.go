package helper

import (
	"strings"

	"ptm-indonesia/config"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Localizer struct {
	bundle           *i18n.Bundle
	matcher          language.Matcher
	defaultLanguage  string
	fallbackLanguage string
}

func NewLocalizer(bundle *i18n.Bundle, cfg *config.AppConfig) *Localizer {
	supportedLanguages := []language.Tag{language.Indonesian, language.English}

	return &Localizer{
		bundle:           bundle,
		matcher:          language.NewMatcher(supportedLanguages),
		defaultLanguage:  cfg.App.DefaultLanguage,
		fallbackLanguage: cfg.App.FallbackLanguage,
	}
}

func (l *Localizer) Resolve(header string) string {
	if strings.TrimSpace(header) == "" {
		return l.defaultLanguage
	}

	tags, _, err := language.ParseAcceptLanguage(header)
	if err != nil || len(tags) == 0 {
		return l.defaultLanguage
	}

	tag, _, _ := l.matcher.Match(tags...)
	base, _ := tag.Base()

	return base.String()
}

func (l *Localizer) MustLocalize(locale string, messageID string) string {
	if strings.TrimSpace(locale) == "" {
		locale = l.defaultLanguage
	}

	localizer := i18n.NewLocalizer(l.bundle, locale, l.fallbackLanguage, l.defaultLanguage)
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return messageID
	}

	return message
}
