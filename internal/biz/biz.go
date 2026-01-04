package biz

import (
	"context"
	"github.com/google/wire"
	"quest-admin/internal/biz/greeter"
	"quest-admin/internal/biz/user"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(greeter.NewGreeterUsecase, user.NewUserUsecase)

type Transaction interface {
	Tx(context.Context, func(ctx context.Context) error) error
}
