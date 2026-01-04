package redis

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

func NewRedSync(rdb *redis.Client) *redsync.Redsync {
	pool := goredis.NewPool(rdb)
	r := redsync.New(pool)
	return r
}
