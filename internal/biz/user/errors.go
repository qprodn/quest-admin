package user

import (
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrUserNotFound            = errors.NotFound("USER_NOT_FOUND", "user not found")
	ErrUserExists              = errors.Conflict("USERNAME_ALREADY_EXISTS", "user already exists")
	ErrInvalidPassword         = errors.BadRequest("INVALID_PASSWORD", "invalid password")
	ErrPasswordConfirmMismatch = errors.BadRequest("PASSWORD_CONFIRM_MISMATCH", "password confirm mismatch")
	ErrInvalidOperationType    = errors.BadRequest("INVALID_OPERATION_TYPE", "invalid operation type")
	ErrPostNotFound            = errors.NotFound("USER_NOT_FOUND", "post not found")
	ErrInternalServer          = errors.InternalServer("UNKNOWN_SYSTEM_ERROR", "UNKNOWN_SYSTEM_ERROR")
	ErrPasswordNotMatch        = errors.BadRequest("PASSWORD_NOT_MATCH", "password not match")
)
