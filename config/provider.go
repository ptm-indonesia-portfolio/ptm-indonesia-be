package config

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAppConfig, NewLogger, NewDatabase, NewI18nBundle)
