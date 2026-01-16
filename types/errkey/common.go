package errkey

import "quest-admin/pkg/errorx"

var (
	ErrInternalServer errorx.ErrorKey = "INTERNAL_SERVER_ERROR"
	ErrBadRequest     errorx.ErrorKey = "BAD_REQUEST"
	ErrUnauthorized   errorx.ErrorKey = "UNAUTHORIZED"
	ErrForbidden      errorx.ErrorKey = "FORBIDDEN"
	ErrNotFound       errorx.ErrorKey = "NOT_FOUND"
	ErrConflict       errorx.ErrorKey = "CONFLICT"
)

func init() {
	errorx.Register(ErrInternalServer, 500, "INTERNAL_SERVER_ERROR", "出错了,请稍后再试")
	errorx.Register(ErrBadRequest, 400, "BAD_REQUEST", "请求参数错误: %s")
	errorx.Register(ErrUnauthorized, 401, "UNAUTHORIZED", "未授权访问")
	errorx.Register(ErrForbidden, 403, "FORBIDDEN", "禁止访问")
	errorx.Register(ErrNotFound, 404, "NOT_FOUND", "资源未找到: %s")
	errorx.Register(ErrConflict, 409, "CONFLICT", "资源冲突: %s")
}
