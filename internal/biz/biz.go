package biz

import (
	"quest-admin/internal/biz/auth"
	"quest-admin/internal/biz/config"
	"quest-admin/internal/biz/dict"
	"quest-admin/internal/biz/organization"
	"quest-admin/internal/biz/permission"
	"quest-admin/internal/biz/tenant"
	"quest-admin/internal/biz/user"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	user.NewUserUsecase,
	organization.NewDepartmentUsecase,
	organization.NewPostUsecase,
	tenant.NewTenantUsecase,
	tenant.NewTenantPackageUsecase,
	permission.NewMenuUsecase,
	permission.NewRoleUsecase,
	config.NewConfigUsecase,
	auth.NewAuthUsecase,
	dict.NewDictTypeUsecase,
	dict.NewDictDataUsecase,
)
