package repositories

import "github.com/google/wire"

var ProviderRepositorySet = wire.NewSet(
	NewRepository,
)
