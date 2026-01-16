package errkey

import "quest-admin/pkg/errorx"

var (
	ErrTenantNotFound      errorx.ErrorKey = "TENANT_NOT_FOUND"
	ErrTenantNameExists    errorx.ErrorKey = "TENANT_NAME_EXISTS"
	ErrTenantHasUsers      errorx.ErrorKey = "TENANT_HAS_USERS"
	ErrInvalidTenantStatus errorx.ErrorKey = "INVALID_TENANT_STATUS"
	ErrInvalidExpireTime   errorx.ErrorKey = "INVALID_TENANT_EXPIRE_TIME"
	ErrInvalidAccountCount errorx.ErrorKey = "INVALID_TENANT_ACCOUNT_COUNT"
)

var (
	ErrTenantPackageNotFound      errorx.ErrorKey = "TENANT_PACKAGE_NOT_FOUND"
	ErrTenantPackageNameExists    errorx.ErrorKey = "TENANT_PACKAGE_NAME_EXISTS"
	ErrTenantPackageInUse         errorx.ErrorKey = "TENANT_PACKAGE_IN_USE"
	ErrInvalidTenantPackageStatus errorx.ErrorKey = "INVALID_TENANT_PACKAGE_STATUS"
)

func init() {
	errorx.Register(ErrTenantPackageNotFound, 404, "TENANT_PACKAGE_NOT_FOUND", "tenant package not found")
	errorx.Register(ErrTenantPackageNameExists, 409, "TENANT_PACKAGE_NAME_EXISTS", "tenant package name already exists")
	errorx.Register(ErrTenantPackageInUse, 400, "TENANT_PACKAGE_IN_USE", "tenant package is in use")
	errorx.Register(ErrInvalidTenantPackageStatus, 400, "INVALID_TENANT_PACKAGE_STATUS", "invalid tenant package status")

	errorx.Register(ErrTenantNotFound, 404, "TENANT_NOT_FOUND", "tenant not found")
	errorx.Register(ErrTenantNameExists, 409, "TENANT_NAME_EXISTS", "tenant name already exists")
	errorx.Register(ErrTenantHasUsers, 400, "TENANT_HAS_USERS", "tenant has users")
	errorx.Register(ErrInvalidTenantStatus, 400, "INVALID_TENANT_STATUS", "invalid tenant status")
	errorx.Register(ErrInvalidExpireTime, 400, "INVALID_TENANT_EXPIRE_TIME", "invalid tenant expire time")
	errorx.Register(ErrInvalidAccountCount, 400, "INVALID_TENANT_ACCOUNT_COUNT", "invalid tenant account count")
}
