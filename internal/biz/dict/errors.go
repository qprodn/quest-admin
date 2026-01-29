package dict

import (
	"quest-admin/pkg/errorx"
	"quest-admin/types/errkey"
)

var (
	ErrDictTypeNotFound   = errorx.Err(errkey.ErrDictTypeNotFound)
	ErrDictTypeCodeExists = errorx.Err(errkey.ErrDictTypeCodeExists)
	ErrDictTypeHasData    = errorx.Err(errkey.ErrDictTypeHasData)

	ErrDictDataNotFound    = errorx.Err(errkey.ErrDictDataNotFound)
	ErrDictDataValueExists = errorx.Err(errkey.ErrDictDataValueExists)
)
