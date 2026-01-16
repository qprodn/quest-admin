package permission

import (
	"quest-admin/pkg/errorx"
	"quest-admin/types/errkey"
)

var (
	ErrRoleNotFound      = errorx.Err(errkey.ErrRoleNotFound)
	ErrRoleNameExists    = errorx.Err(errkey.ErrRoleNameExists)
	ErrRoleCodeExists    = errorx.Err(errkey.ErrRoleCodeExists)
	ErrRoleHasUsers      = errorx.Err(errkey.ErrRoleHasUsers)
	ErrInvalidRoleStatus = errorx.Err(errkey.ErrInvalidRoleStatus)
)

var (
	ErrMenuNotFound      = errorx.Err(errkey.ErrMenuNotFound)
	ErrMenuNameExists    = errorx.Err(errkey.ErrMenuNameExists)
	ErrMenuHasChildren   = errorx.Err(errkey.ErrMenuHasChildren)
	ErrInvalidParentMenu = errorx.Err(errkey.ErrInvalidParentMenu)
	ErrInvalidMenuStatus = errorx.Err(errkey.ErrInvalidMenuStatus)
	ErrInvalidMenuType   = errorx.Err(errkey.ErrInvalidMenuType)
	ErrMenuLevelExceeded = errorx.Err(errkey.ErrMenuLevelExceeded)
	ErrInvalidMenuPath   = errorx.Err(errkey.ErrInvalidMenuPath)
)

var (
	ErrInternalServer = errorx.Err(errkey.ErrInternalServer)
)
