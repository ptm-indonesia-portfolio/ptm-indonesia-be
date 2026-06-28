package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	App       AppSection
	Admin     AdminSection
	Auth      AuthSection
	RateLimit RateLimitSection
	Database  DatabaseSection
	Log       LogSection
	Migration MigrationSection
}

type AppSection struct {
	Name              string
	Environment       string
	Host              string
	Port              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	Timezone          string
	DefaultLanguage   string
	FallbackLanguage  string
	SupportedLanguage []string
	CORSOrigins       []string
}

type AdminSection struct {
	Email string
}

type AuthSection struct {
	BaseURL            string
	FrontendURL        string
	CookieSecret       string
	CookieName         string
	RefreshCookieName  string
	LoggedInCookieName string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	GoogleIssuerURL    string
}

type RateLimitSection struct {
	MaxRequests int
	Window      time.Duration
}

type DatabaseSection struct {
	Host               string
	Port               int
	Name               string
	User               string
	Password           string
	SSLMode            string
	Timezone           string
	MaxIdleConns       int
	MaxOpenConns       int
	ConnMaxLifetime    time.Duration
	HealthCheckTimeout time.Duration
}

type LogSection struct {
	FilePath string
	Level    string
}

type MigrationSection struct {
	Source string
}

func NewAppConfig() (*AppConfig, error) {
	v := viper.New()
	configFilePath, hasExplicitConfigPath := os.LookupEnv("ENV_FILE_PATH")
	configFilePath = strings.TrimSpace(configFilePath)
	if configFilePath == "" {
		hasExplicitConfigPath = false
		configFilePath = ".env"
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(v)

	if _, err := os.Stat(configFilePath); err == nil {
		v.SetConfigFile(configFilePath)
		v.SetConfigType("env")

		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config file: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("stat config file: %w", err)
	} else if hasExplicitConfigPath {
		return nil, fmt.Errorf("config file not found: %s", configFilePath)
	}

	cfg := &AppConfig{
		App: AppSection{
			Name:              v.GetString("APP_NAME"),
			Environment:       v.GetString("APP_ENV"),
			Host:              v.GetString("APP_HOST"),
			Port:              v.GetString("APP_PORT"),
			ReadTimeout:       time.Duration(v.GetInt("APP_READ_TIMEOUT_SECONDS")) * time.Second,
			WriteTimeout:      time.Duration(v.GetInt("APP_WRITE_TIMEOUT_SECONDS")) * time.Second,
			IdleTimeout:       time.Duration(v.GetInt("APP_IDLE_TIMEOUT_SECONDS")) * time.Second,
			ShutdownTimeout:   time.Duration(v.GetInt("APP_SHUTDOWN_TIMEOUT_SECONDS")) * time.Second,
			Timezone:          v.GetString("APP_TIMEZONE"),
			DefaultLanguage:   strings.ToLower(v.GetString("APP_DEFAULT_LANGUAGE")),
			FallbackLanguage:  strings.ToLower(v.GetString("APP_FALLBACK_LANGUAGE")),
			SupportedLanguage: []string{"id", "en"},
			CORSOrigins:       buildCORSOrigins(v.GetString("CORS_ORIGINS"), v.GetString("FRONTEND_URL")),
		},
		Admin: AdminSection{
			Email: strings.TrimSpace(v.GetString("EMAIL_ADMIN")),
		},
		Auth: AuthSection{
			BaseURL:            strings.TrimSpace(v.GetString("BASEURL")),
			FrontendURL:        strings.TrimSpace(v.GetString("FRONTEND_URL")),
			CookieSecret:       v.GetString("COOKIE_SECRET"),
			CookieName:         v.GetString("AUTH_COOKIE_NAME"),
			RefreshCookieName:  v.GetString("AUTH_REFRESH_COOKIE_NAME"),
			LoggedInCookieName: v.GetString("AUTH_LOGGED_IN_COOKIE_NAME"),
			AccessTokenTTL:     time.Duration(v.GetInt("AUTH_ACCESS_TOKEN_TTL_MINUTES")) * time.Minute,
			RefreshTokenTTL:    time.Duration(v.GetInt("AUTH_REFRESH_TOKEN_TTL_HOURS")) * time.Hour,
			GoogleClientID:     strings.TrimSpace(v.GetString("GOOGLE_CLIENT_ID")),
			GoogleClientSecret: strings.TrimSpace(v.GetString("GOOGLE_CLIENT_SECRET")),
			GoogleRedirectURL:  strings.TrimSpace(v.GetString("GOOGLE_REDIRECT_URI")),
			GoogleIssuerURL:    strings.TrimSpace(v.GetString("GOOGLE_ISSUER_URL")),
		},
		RateLimit: RateLimitSection{
			MaxRequests: v.GetInt("RATE_LIMIT_MAX_REQUESTS"),
			Window:      time.Duration(v.GetInt("RATE_LIMIT_WINDOW_MS")) * time.Millisecond,
		},
		Database: DatabaseSection{
			Host:               v.GetString("DB_HOST"),
			Port:               v.GetInt("DB_PORT"),
			Name:               v.GetString("DB_NAME"),
			User:               v.GetString("DB_USER"),
			Password:           v.GetString("DB_PASSWORD"),
			SSLMode:            v.GetString("DB_SSLMODE"),
			Timezone:           v.GetString("DB_TIMEZONE"),
			MaxIdleConns:       v.GetInt("DB_MAX_IDLE_CONNS"),
			MaxOpenConns:       v.GetInt("DB_MAX_OPEN_CONNS"),
			ConnMaxLifetime:    time.Duration(v.GetInt("DB_CONN_MAX_LIFETIME_MINUTES")) * time.Minute,
			HealthCheckTimeout: time.Duration(v.GetInt("DB_HEALTH_CHECK_TIMEOUT_SECONDS")) * time.Second,
		},
		Log: LogSection{
			FilePath: v.GetString("LOG_FILE_PATH"),
			Level:    strings.ToLower(v.GetString("LOG_LEVEL")),
		},
		Migration: MigrationSection{
			Source: v.GetString("MIGRATION_SOURCE"),
		},
	}

	if strings.TrimSpace(cfg.Admin.Email) == "" {
		return nil, fmt.Errorf("EMAIL_ADMIN is required")
	}

	return cfg, nil
}

func (c *AppConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.App.Host, c.App.Port)
}

func (c *AppConfig) DatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
		c.Database.Timezone,
	)
}

func (c *AppConfig) MigrationDatabaseURL() string {
	query := url.Values{}
	query.Set("sslmode", c.Database.SSLMode)

	if c.Database.Timezone != "" {
		query.Set("timezone", c.Database.Timezone)
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?%s",
		url.QueryEscape(c.Database.User),
		url.QueryEscape(c.Database.Password),
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		query.Encode(),
	)
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_NAME", "PTM Indonesia API")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_HOST", "0.0.0.0")
	v.SetDefault("APP_PORT", "9100")
	v.SetDefault("APP_READ_TIMEOUT_SECONDS", 10)
	v.SetDefault("APP_WRITE_TIMEOUT_SECONDS", 10)
	v.SetDefault("APP_IDLE_TIMEOUT_SECONDS", 60)
	v.SetDefault("APP_SHUTDOWN_TIMEOUT_SECONDS", 10)
	v.SetDefault("APP_TIMEZONE", "Asia/Jakarta")
	v.SetDefault("APP_DEFAULT_LANGUAGE", "id")
	v.SetDefault("APP_FALLBACK_LANGUAGE", "en")
	v.SetDefault("EMAIL_ADMIN", "")
	v.SetDefault("CORS_ORIGINS", "http://localhost:3100,http://localhost:3101")
	v.SetDefault("COOKIE_SECRET", "change_me")
	v.SetDefault("AUTH_COOKIE_NAME", "ptm_auth_token")
	v.SetDefault("AUTH_REFRESH_COOKIE_NAME", "ptm_refresh_token")
	v.SetDefault("AUTH_LOGGED_IN_COOKIE_NAME", "logged_in")
	v.SetDefault("AUTH_ACCESS_TOKEN_TTL_MINUTES", 15)
	v.SetDefault("AUTH_REFRESH_TOKEN_TTL_HOURS", 8760)
	v.SetDefault("GOOGLE_REDIRECT_URI", "http://localhost:9100/api/v1/auth/google/callback")
	v.SetDefault("GOOGLE_ISSUER_URL", "https://accounts.google.com")
	v.SetDefault("BASEURL", "http://localhost:9100")
	v.SetDefault("FRONTEND_URL", "http://localhost:3101/home")
	v.SetDefault("RATE_LIMIT_MAX_REQUESTS", 60)
	v.SetDefault("RATE_LIMIT_WINDOW_MS", 60000)

	v.SetDefault("DB_HOST", "db")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_NAME", "ptm_indonesia")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "postgres")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("DB_TIMEZONE", "Asia/Jakarta")
	v.SetDefault("DB_MAX_IDLE_CONNS", 10)
	v.SetDefault("DB_MAX_OPEN_CONNS", 25)
	v.SetDefault("DB_CONN_MAX_LIFETIME_MINUTES", 30)
	v.SetDefault("DB_HEALTH_CHECK_TIMEOUT_SECONDS", 3)

	v.SetDefault("LOG_FILE_PATH", "logs/app.log")
	v.SetDefault("LOG_LEVEL", "error")

	v.SetDefault("MIGRATION_SOURCE", "file://db/migrations")
}

func parseCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{}
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		result = append(result, trimmed)
	}

	return result
}

func buildCORSOrigins(origins string, frontendURL string) []string {
	result := parseCSV(origins)
	frontendOrigin := extractURLOrigin(frontendURL)
	if frontendOrigin == "" || containsString(result, frontendOrigin) {
		return result
	}

	return append(result, frontendOrigin)
}

func extractURLOrigin(rawURL string) string {
	trimmedURL := strings.TrimSpace(rawURL)
	if trimmedURL == "" {
		return ""
	}

	parsedURL, err := url.Parse(trimmedURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return ""
	}

	return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}
