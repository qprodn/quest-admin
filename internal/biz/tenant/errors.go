package tenant

import (
	v1 "quest-admin/api/gen/tenant/v1"

	"github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrTenantPackageNotFound      = errors.NotFound("TENANT_PACKAGE_NOT_FOUND", "tenant package not found")
	ErrTenantPackageNameExists    = errors.Conflict("TENANT_PACKAGE_NAME_EXISTS", "tenant package name already exists")
	ErrTenantPackageInUse         = errors.BadRequest("TENANT_PACKAGE_IN_USE", "tenant package is in use")
	ErrInvalidTenantPackageStatus = errors.BadRequest("INVALID_TENANT_PACKAGE_STATUS", "invalid tenant package status")
)

var (
	ErrTenantNotFound      = errors.NotFound("TENANT_NOT_FOUND", "tenant not found")
	ErrTenantNameExists    = errors.Conflict("TENANT_NAME_EXISTS", "tenant name already exists")
	ErrTenantHasUsers      = errors.BadRequest("TENANT_HAS_USERS", "tenant has users")
	ErrInvalidTenantStatus = errors.BadRequest("INVALID_TENANT_STATUS", "invalid tenant status")
	ErrInvalidExpireTime   = errors.BadRequest("INVALID_TENANT_EXPIRE_TIME", "invalid tenant expire time")
	ErrInvalidAccountCount = errors.BadRequest("INVALID_TENANT_ACCOUNT_COUNT", "invalid tenant account count")
)
