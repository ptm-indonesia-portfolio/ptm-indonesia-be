//go:build wireinject

package bootstrap

import (
	"ptm-indonesia/config"
	"ptm-indonesia/controllers"
	"ptm-indonesia/helper"
	"ptm-indonesia/repository"
	"ptm-indonesia/services"
	"ptm-indonesia/validation"

	"github.com/google/wire"
)

func InitializeHTTPApplication() (*HTTPApplication, func(), error) {
	wire.Build(
		config.ProviderSet,
		helper.ProviderSet,
		validation.ProviderSet,
		repository.ProviderSet,
		services.ProviderSet,
		controllers.ProviderSet,
		ProviderSet,
	)

	return nil, nil, nil
}
