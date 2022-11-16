package services

import "github.com/google/wire"

var ProviderServiceSet = wire.NewSet(
	NewService,
)
