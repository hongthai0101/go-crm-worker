package clients

import "github.com/google/wire"

var ProviderHttpClientSet = wire.NewSet(
	NewHttpClient,
)
