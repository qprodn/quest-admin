package data

import (
	"quest-admin/internal/data/auth"
	"quest-admin/internal/data/config"
	"quest-admin/internal/data/data"
	"quest-admin/internal/data/idgen"
	"quest-admin/internal/data/organization"
	"quest-admin/internal/data/permission"
	"quest-admin/internal/data/pg"
	"quest-admin/internal/data/redis"
	"quest-admin/internal/data/tenant"
	"quest-admin/internal/data/transaction"
	"quest-admin/internal/data/user"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	data.NewData,
	transaction.NewManager,
	idgen.NewIDGenerator,
	redis.NewRedis,
	redis.NewRedSync,
	pg.NewDB,
	user.NewUserRepo,
	user.NewUserRoleRepo,
	user.NewUserPostRepo,
	user.NewUserDeptRepo,
	organization.NewDepartmentRepo,
	organization.NewPostRepo,
	tenant.NewTenantRepo,
	tenant.NewTenantPackageRepo,
	permission.NewRoleRepo,
	permission.NewMenuRepo,
	permission.NewRoleMenuRepo,
	config.NewConfigRepo,
	auth.NewAuthManager,
)
