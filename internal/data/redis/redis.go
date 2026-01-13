package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"quest-admin/internal/conf"
	"time"
)

func NewRedis(c *conf.Bootstrap) *redis.Client {
	// 创建Redis连接
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Data.Redis.Addr,
		Password: c.Data.Redis.Password,
		DB:       int(c.Data.Redis.Db),

		// 超时设置
		ReadTimeout:  time.Duration(c.Data.Redis.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.Data.Redis.ReadTimeout) * time.Second,

		// 从配置文件中获取连接池配置
		PoolSize:     int(c.Data.Redis.PoolSize),     // 连接池大小
		MinIdleConns: int(c.Data.Redis.MinIdleConns), // 最小空闲连接数
	})

	// 测试Redis连接
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}
	return rdb
}
