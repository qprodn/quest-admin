package transaction

import (
	"context"

	"github.com/uptrace/bun"
)

type Manager interface {
	Tx(ctx context.Context, fn func(ctx context.Context) error) error
}

type ContextTxKey struct{}

type manager struct {
	db *bun.DB
}

func NewManager(db *bun.DB) Manager {
	return &manager{db: db}
}
func (m *manager) Tx(ctx context.Context, fn func(ctx context.Context) error) error {
	err := m.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		ctx = context.WithValue(ctx, ContextTxKey{}, tx)
		return fn(ctx)
	})
	return err
}
