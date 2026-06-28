package services

import (
	contract "ptm-indonesia/services/contract"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewHealthService,
	NewAuthService,
	NewUserService,
	wire.Bind(new(contract.HealthService), new(*HealthService)),
	wire.Bind(new(contract.AuthService), new(*AuthService)),
	wire.Bind(new(contract.UserService), new(*UserService)),
)
