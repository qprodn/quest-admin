package auth

import (
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	// ErrUserNotFound 用户未找到
	ErrUserNotFound = errors.NotFound("USER_NOT_FOUND", "用户未找到")

	// ErrPasswordError 密码错误
	ErrPasswordError = errors.Unauthorized("PASSWORD_ERROR", "密码错误")

	// ErrUserDisabled 用户已被禁用
	ErrUserDisabled = errors.Forbidden("USER_DISABLED", "用户已被禁用")

	// ErrTokenInvalid Token 无效
	ErrTokenInvalid = errors.Unauthorized("TOKEN_INVALID", "Token 无效")

	// ErrTokenExpired Token 过期
	ErrTokenExpired = errors.Unauthorized("TOKEN_EXPIRED", "Token 过期")

	// ErrInternalServer 内部服务器错误
	ErrInternalServer = errors.InternalServer("INTERNAL_SERVER_ERROR", "内部服务器错误")
)
