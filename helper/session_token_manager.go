package helper

import (
	"fmt"
	"time"

	"ptm-indonesia/config"
	"ptm-indonesia/model"

	"github.com/golang-jwt/jwt/v5"
)

type SessionTokenManager struct {
	cookieSecret string
	issuer       string
	sessionTTL   time.Duration
}

type sessionClaims struct {
	jwt.RegisteredClaims
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	GoogleID *string `json:"google_id,omitempty"`
	Status   int16   `json:"status"`
}

func NewSessionTokenManager(cfg *config.AppConfig) *SessionTokenManager {
	return &SessionTokenManager{
		cookieSecret: cfg.Auth.CookieSecret,
		issuer:       cfg.App.Name,
		sessionTTL:   cfg.Auth.AccessTokenTTL,
	}
}

func (m *SessionTokenManager) Generate(user *model.User) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, sessionClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.sessionTTL)),
		},
		Name:     user.Name,
		Email:    user.Email,
		GoogleID: user.GoogleID,
		Status:   user.Status,
	})

	return token.SignedString([]byte(m.cookieSecret))
}

func (m *SessionTokenManager) Parse(tokenString string) (*model.AuthSessionUser, error) {
	claims := &sessionClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(m.cookieSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer(m.issuer))
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid auth token")
	}

	var userID uint64
	if _, err := fmt.Sscanf(claims.Subject, "%d", &userID); err != nil {
		return nil, fmt.Errorf("parse auth token subject: %w", err)
	}

	return &model.AuthSessionUser{
		ID:       userID,
		Name:     claims.Name,
		Email:    claims.Email,
		GoogleID: claims.GoogleID,
		Status:   claims.Status,
	}, nil
}
