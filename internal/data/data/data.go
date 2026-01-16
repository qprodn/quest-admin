package data

import (
	"context"

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

type contextTxKey struct{}

func NewData(db *bun.DB) *Data {
	return &Data{
		Db: db,
	}
}

func (d *Data) DB(ctx context.Context) bun.IDB {
	idb, ok := ctx.Value(contextTxKey{}).(bun.IDB)
	if ok {
		return idb
	}
	return d.Db
}

func (d *Data) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.Db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}
