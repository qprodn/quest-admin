package config

import (
	"quest-admin/pkg/errorx"
	"quest-admin/types/errkey"
)

var (
	ErrConfigNotFound = errorx.Err(errkey.ErrConfigNotFound)
	ErrConfigExists   = errorx.Err(errkey.ErrConfigExists)
	ErrConfigDisabled = errorx.Err(errkey.ErrConfigDisabled)
)
