package auth

import (
	"github.com/click33/sa-token-go/core"
	storage "github.com/click33/sa-token-go/storage/redis"
	"github.com/click33/sa-token-go/stputil"
	"github.com/redis/go-redis/v9"
)

type Manager struct {
	Admin *stputil.StpLogic
}

type Admin stputil.StpLogic

func NewAuthManager(redisClient *redis.Client) *Manager {
	admin := stputil.NewStpLogic(
		core.NewBuilder().
			Storage(storage.NewStorageFromClient(redisClient)).
			KeyPrefix("qa:admin:").
			TokenName("Authorization").
			Timeout(43200).
			TokenStyle(core.TokenStyleTik).
			IsPrintBanner(false).
			Build())
	return &Manager{Admin: admin}
}
