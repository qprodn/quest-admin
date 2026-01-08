package biz

import (
	"context"
	"github.com/google/wire"
	"quest-admin/internal/biz/greeter"
	"quest-admin/internal/biz/organization"
	"quest-admin/internal/biz/tenant"
	"quest-admin/internal/biz/user"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	greeter.NewGreeterUsecase,
	user.NewUserUsecase,
	user.NewUserRoleUsecase,
	organization.NewDepartmentUsecase,
	organization.NewPostUsecase,
	tenant.NewTenantUsecase,
	tenant.NewTenantPackageUsecase,
)

type Transaction interface {
	Tx(context.Context, func(ctx context.Context) error) error
}
