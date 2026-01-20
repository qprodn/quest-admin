package data

import (
	"context"
	"quest-admin/internal/data/transaction"

	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

// Data .
type Data struct {
	Db    bun.IDB
	Rdb   *redis.Client
	Rsync *redsync.Redsync
}

func NewData(db *bun.DB) *Data {
	return &Data{
		Db: db,
	}
}

func (d *Data) DB(ctx context.Context) bun.IDB {
	idb, ok := ctx.Value(transaction.ContextTxKey{}).(bun.IDB)
	if ok {
		return idb
	}
	return d.Db
}
