package helper

import (
	"context"
	"fmt"

	"ptm-indonesia/config"
	"ptm-indonesia/model"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type GoogleOIDCClient struct {
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
}

func NewGoogleOIDCClient(cfg *config.AppConfig) (*GoogleOIDCClient, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, cfg.Auth.GoogleIssuerURL)
	if err != nil {
		return nil, fmt.Errorf("discover google oidc provider: %w", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID:     cfg.Auth.GoogleClientID,
		ClientSecret: cfg.Auth.GoogleClientSecret,
		RedirectURL:  cfg.Auth.GoogleRedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	return &GoogleOIDCClient{
		oauth2Config: oauth2Config,
		verifier: provider.Verifier(&oidc.Config{
			ClientID: cfg.Auth.GoogleClientID,
		}),
	}, nil
}

func (g *GoogleOIDCClient) AuthCodeURL(state string) string {
	return g.oauth2Config.AuthCodeURL(state)
}

func (g *GoogleOIDCClient) ExchangeCode(ctx context.Context, code string) (*model.GoogleIdentity, error) {
	oauth2Token, err := g.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange google auth code: %w", err)
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("missing google id token")
	}

	idToken, err := g.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("verify google id token: %w", err)
	}

	var identity model.GoogleIdentity
	if err := idToken.Claims(&identity); err != nil {
		return nil, fmt.Errorf("extract google identity claims: %w", err)
	}

	return &identity, nil
}
