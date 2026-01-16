package user

import (
	"quest-admin/pkg/errorx"
	"quest-admin/types/errkey"
)

var (
	ErrUserNotFound            = errorx.Err(errkey.ErrUserNotFound)
	ErrUserExists              = errorx.Err(errkey.ErrUserExists)
	ErrInvalidPassword         = errorx.Err(errkey.ErrInvalidPassword)
	ErrPasswordConfirmMismatch = errorx.Err(errkey.ErrPasswordConfirmMismatch)
	ErrInvalidOperationType    = errorx.Err(errkey.ErrInvalidOperationType)
	ErrInternalServer          = errorx.Err(errkey.ErrInternalServer)
	ErrPasswordNotMatch        = errorx.Err(errkey.ErrPasswordNotMatch)
	ErrUserDisabled            = errorx.Err(errkey.ErrUserDisabled)
)
