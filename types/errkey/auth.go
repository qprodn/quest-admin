package errkey

import "quest-admin/pkg/errorx"

var (
	ErrTokenInvalid  errorx.ErrorKey = "TOKEN_INVALID"
	ErrTokenExpired  errorx.ErrorKey = "TOKEN_EXPIRED"
	ErrPasswordError errorx.ErrorKey = "PASSWORD_ERROR"
)

func init() {
	errorx.Register(ErrTokenInvalid, 401, "TOKEN_INVALID", "token invalid")
	errorx.Register(ErrTokenExpired, 401, "TOKEN_EXPIRED", "token expired")
	errorx.Register(ErrPasswordError, 401, "PASSWORD_ERROR", "password error")
}
