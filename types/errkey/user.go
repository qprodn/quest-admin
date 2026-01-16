package errkey

import "quest-admin/pkg/errorx"

var (
	ErrUserNotFound            errorx.ErrorKey = "USER_NOT_FOUND"
	ErrUserExists              errorx.ErrorKey = "USER_EXISTS"
	ErrInvalidPassword         errorx.ErrorKey = "INVALID_PASSWORD"
	ErrPasswordConfirmMismatch errorx.ErrorKey = "PASSWORD_CONFIRM_MISMATCH"
	ErrInvalidOperationType    errorx.ErrorKey = "INVALID_OPERATION_TYPE"
	ErrPasswordNotMatch        errorx.ErrorKey = "PASSWORD_NOT_MATCH"
	ErrUserDisabled            errorx.ErrorKey = "USER_DISABLED"
)

func init() {
	errorx.Register(ErrUserNotFound, 404, "USER_NOT_FOUND", "user not found")
	errorx.Register(ErrUserExists, 409, "USERNAME_ALREADY_EXISTS", "user already exists")
	errorx.Register(ErrInvalidPassword, 400, "INVALID_PASSWORD", "invalid password")
	errorx.Register(ErrPasswordConfirmMismatch, 400, "PASSWORD_CONFIRM_MISMATCH", "password confirm mismatch")
	errorx.Register(ErrInvalidOperationType, 400, "INVALID_OPERATION_TYPE", "invalid operation type")
	errorx.Register(ErrPasswordNotMatch, 400, "PASSWORD_NOT_MATCH", "password not match")
	errorx.Register(ErrUserDisabled, 403, "USER_DISABLED", "user disabled")
}
