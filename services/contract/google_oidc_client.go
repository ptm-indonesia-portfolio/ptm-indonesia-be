package contract

import (
	"context"

	"ptm-indonesia/model"
)

type GoogleOIDCClient interface {
	AuthCodeURL(state string) string
	ExchangeCode(ctx context.Context, code string) (*model.GoogleIdentity, error)
}
