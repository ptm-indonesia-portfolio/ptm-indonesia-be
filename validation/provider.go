package validation

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewRequestValidator)
