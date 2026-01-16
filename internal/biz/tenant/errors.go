package tenant

import (
	"quest-admin/pkg/errorx"
	"quest-admin/types/errkey"
)

var (
	ErrTenantPackageNotFound      = errorx.Err(errkey.ErrTenantPackageNotFound)
	ErrTenantPackageNameExists    = errorx.Err(errkey.ErrTenantPackageNameExists)
	ErrTenantPackageInUse         = errorx.Err(errkey.ErrTenantPackageInUse)
	ErrInvalidTenantPackageStatus = errorx.Err(errkey.ErrInvalidTenantPackageStatus)
)

var (
	ErrTenantNotFound      = errorx.Err(errkey.ErrTenantNotFound)
	ErrTenantNameExists    = errorx.Err(errkey.ErrTenantNameExists)
	ErrTenantHasUsers      = errorx.Err(errkey.ErrTenantHasUsers)
	ErrInvalidTenantStatus = errorx.Err(errkey.ErrInvalidTenantStatus)
	ErrInvalidExpireTime   = errorx.Err(errkey.ErrInvalidExpireTime)
	ErrInvalidAccountCount = errorx.Err(errkey.ErrInvalidAccountCount)
)
