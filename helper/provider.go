package helper

import (
	servicesContract "ptm-indonesia/services/contract"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewResponder,
	NewLocalizer,
	NewGoogleOIDCClient,
	NewSessionTokenManager,
	NewRefreshTokenManager,
	wire.Bind(new(servicesContract.GoogleOIDCClient), new(*GoogleOIDCClient)),
	wire.Bind(new(servicesContract.SessionTokenManager), new(*SessionTokenManager)),
	wire.Bind(new(servicesContract.RefreshTokenManager), new(*RefreshTokenManager)),
)
