package organization

import (
	v1 "quest-admin/api/gen/organization/v1"

	"github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrDepartmentNotFound      = errors.NotFound("DEPARTMENT_NOT_FOUND", "department not found")
	ErrDepartmentNameExists    = errors.Conflict("DEPARTMENT_NAME_EXISTS", "department name already exists")
	ErrDepartmentHasChildren   = errors.BadRequest("DEPARTMENT_HAS_CHILDREN", "department has children")
	ErrDepartmentHasUsers      = errors.BadRequest("DEPARTMENT_HAS_USERS", "department has users")
	ErrInvalidParentDepartment = errors.BadRequest("INVALID_PARENT_DEPARTMENT", "invalid parent department")
)

var (
	ErrPostNotFound   = errors.NotFound("POST_NOT_FOUND", "post not found")
	ErrPostNameExists = errors.Conflict("POST_NAME_EXISTS", "post name already exists")
	ErrPostHasUsers   = errors.BadRequest(
		"POST_HAS_USERS", "post has users")
)
