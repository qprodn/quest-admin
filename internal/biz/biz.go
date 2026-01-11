package biz

import (
	"context"
	"quest-admin/internal/biz/greeter"
	"quest-admin/internal/biz/organization"
	"quest-admin/internal/biz/permission"
	"quest-admin/internal/biz/tenant"
	"quest-admin/internal/biz/user"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	greeter.NewGreeterUsecase,
	user.NewUserUsecase,
	user.NewUserRoleUsecase,
	user.NewUserPostUsecase,
	user.NewUserDeptUsecase,
	organization.NewDepartmentUsecase,
	organization.NewPostUsecase,
	tenant.NewTenantUsecase,
	tenant.NewTenantPackageUsecase,
	permission.NewMenuUsecase,
	permission.NewRoleUsecase,
)

type Transaction interface {
	Tx(context.Context, func(ctx context.Context) error) error
}
