package service

import (
	"github.com/google/wire"
	"quest-admin/internal/service/greeter"
	"quest-admin/internal/service/user"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(greeter.NewGreeterService, user.NewUserService)
