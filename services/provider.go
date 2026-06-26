package services

import (
	contract "ptm-indonesia/services/contract"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewHealthService,
	NewAuthService,
	wire.Bind(new(contract.HealthService), new(*HealthService)),
	wire.Bind(new(contract.AuthService), new(*AuthService)),
)
