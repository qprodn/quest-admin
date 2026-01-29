package errkey

import "quest-admin/pkg/errorx"

var (
	ErrConfigNotFound errorx.ErrorKey = "CONFIG_NOT_FOUND"
	ErrConfigExists   errorx.ErrorKey = "CONFIG_EXISTS"
	ErrConfigDisabled errorx.ErrorKey = "CONFIG_DISABLED"
)

func init() {
	errorx.Register(ErrConfigNotFound, 404, "CONFIG_NOT_FOUND", "配置不存在")
	errorx.Register(ErrConfigExists, 409, "CONFIG_EXISTS", "配置已存在")
	errorx.Register(ErrConfigDisabled, 403, "CONFIG_DISABLED", "配置已禁用")
}
