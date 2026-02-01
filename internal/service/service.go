package service

import (
	"quest-admin/internal/service/auth"
	"quest-admin/internal/service/config"
	"quest-admin/internal/service/dict"
	"quest-admin/internal/service/organization"
	"quest-admin/internal/service/permission"
	"quest-admin/internal/service/tenant"
	"quest-admin/internal/service/user"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	user.NewUserService,
	tenant.NewTenantService,
	tenant.NewTenantPackageService,
	permission.NewMenuService,
	permission.NewRoleService,
	organization.NewDepartmentService,
	organization.NewPostService,
	config.NewConfigService,
	auth.NewAuthService,
	dict.NewDictService,
)
