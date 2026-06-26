package repository

import (
	contract "ptm-indonesia/repository/contract"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewSystemRepository,
	NewUserRepository,
	NewRefreshTokenRepository,
	wire.Bind(new(contract.SystemRepository), new(*SystemRepository)),
	wire.Bind(new(contract.UserRepository), new(*UserRepository)),
	wire.Bind(new(contract.RefreshTokenRepository), new(*RefreshTokenRepository)),
)
