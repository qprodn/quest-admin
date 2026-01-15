package service

import (
	"quest-admin/internal/service/auth"
	"quest-admin/internal/service/greeter"
	"quest-admin/internal/service/organization"
	"quest-admin/internal/service/permission"
	"quest-admin/internal/service/tenant"
	"quest-admin/internal/service/user"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	greeter.NewGreeterService,
	user.NewUserService,
	user.NewUserRoleService,
	tenant.NewTenantService,
	permission.NewMenuService,
	permission.NewRoleService,
	organization.NewDepartmentService,
	organization.NewPostService,
	auth.NewAuthService,
)
