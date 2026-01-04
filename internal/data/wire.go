package data

import (
	"github.com/google/wire"
	"quest-admin/internal/data/data"
	"quest-admin/internal/data/greeter"
	"quest-admin/internal/data/pg"
	"quest-admin/internal/data/redis"
	"quest-admin/internal/data/user"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	data.NewData,
	redis.NewRedis,
	redis.NewRedSync,
	pg.NewDB,
	greeter.NewGreeterRepo,
	user.NewUserRepo,
)
