package model

type GoogleIdentity struct {
	Subject       string  `json:"sub"`
	Email         string  `json:"email"`
	EmailVerified bool    `json:"email_verified"`
	Name          string  `json:"name"`
	Picture       *string `json:"picture"`
}
