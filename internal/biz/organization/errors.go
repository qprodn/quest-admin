package organization

import (
	v1 "quest-admin/api/gen/organization/v1"

	"github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrDepartmentNotFound      = errors.NotFound(v1.ErrorReason_DEPARTMENT_NOT_FOUND.String(), "department not found")
	ErrDepartmentNameExists    = errors.Conflict(v1.ErrorReason_DEPARTMENT_NAME_EXISTS.String(), "department name already exists")
	ErrDepartmentHasChildren   = errors.BadRequest(v1.ErrorReason_DEPARTMENT_HAS_CHILDREN.String(), "department has children")
	ErrDepartmentHasUsers      = errors.BadRequest(v1.ErrorReason_DEPARTMENT_HAS_USERS.String(), "department has users")
	ErrInvalidParentDepartment = errors.BadRequest(v1.ErrorReason_INVALID_PARENT_DEPARTMENT.String(), "invalid parent department")
)

var (
	ErrPostNotFound   = errors.NotFound(v1.ErrorReason_POST_NOT_FOUND.String(), "post not found")
	ErrPostNameExists = errors.Conflict(v1.ErrorReason_POST_NAME_EXISTS.String(), "post name already exists")
	ErrPostHasUsers   = errors.BadRequest(v1.ErrorReason_POST_HAS_USERS.String(), "post has users")
)
