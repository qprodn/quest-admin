package tenant

import (
	v1 "quest-admin/api/gen/tenant/v1"

	"github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrTenantPackageNotFound      = errors.NotFound(v1.ErrorReason_TENANT_PACKAGE_NOT_FOUND.String(), "tenant package not found")
	ErrTenantPackageNameExists    = errors.Conflict(v1.ErrorReason_TENANT_PACKAGE_NAME_EXISTS.String(), "tenant package name already exists")
	ErrTenantPackageInUse         = errors.BadRequest(v1.ErrorReason_TENANT_PACKAGE_IN_USE.String(), "tenant package is in use")
	ErrInvalidTenantPackageStatus = errors.BadRequest(v1.ErrorReason_INVALID_TENANT_PACKAGE_STATUS.String(), "invalid tenant package status")
)

var (
	ErrTenantNotFound      = errors.NotFound(v1.ErrorReason_TENANT_NOT_FOUND.String(), "tenant not found")
	ErrTenantNameExists    = errors.Conflict(v1.ErrorReason_TENANT_NAME_EXISTS.String(), "tenant name already exists")
	ErrTenantHasUsers      = errors.BadRequest(v1.ErrorReason_TENANT_HAS_USERS.String(), "tenant has users")
	ErrInvalidTenantStatus = errors.BadRequest(v1.ErrorReason_INVALID_TENANT_STATUS.String(), "invalid tenant status")
	ErrInvalidExpireTime   = errors.BadRequest(v1.ErrorReason_INVALID_TENANT_EXPIRE_TIME.String(), "invalid tenant expire time")
	ErrInvalidAccountCount = errors.BadRequest(v1.ErrorReason_INVALID_TENANT_ACCOUNT_COUNT.String(), "invalid tenant account count")
)
