package model

type AuthLoginResult struct {
	AccessToken  string
	RefreshToken string
	User         *AuthSessionUser
}
