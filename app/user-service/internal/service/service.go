package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewFirebaseService,
	NewGithubService,
	NewJwtService,
	NewUserSessionService,
	NewCacheService,
	NewUserService,
)
