package organization

import (
	"quest-admin/pkg/errorx"
	"quest-admin/types/errkey"
)

var (
	ErrDepartmentNotFound      = errorx.Err(errkey.ErrDepartmentNotFound)
	ErrDepartmentNameExists    = errorx.Err(errkey.ErrDepartmentNameExists)
	ErrDepartmentHasChildren   = errorx.Err(errkey.ErrDepartmentHasChildren)
	ErrDepartmentHasUsers      = errorx.Err(errkey.ErrDepartmentHasUsers)
	ErrInvalidParentDepartment = errorx.Err(errkey.ErrInvalidParentDepartment)
)

var (
	ErrPostNotFound   = errorx.Err(errkey.ErrPostNotFound)
	ErrPostNameExists = errorx.Err(errkey.ErrPostNameExists)
	ErrPostHasUsers   = errorx.Err(errkey.ErrPostHasUsers)
)
