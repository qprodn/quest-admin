package errkey

import "quest-admin/pkg/errorx"

var (
	ErrDepartmentNotFound      errorx.ErrorKey = "DEPARTMENT_NOT_FOUND"
	ErrDepartmentNameExists    errorx.ErrorKey = "DEPARTMENT_NAME_EXISTS"
	ErrDepartmentHasChildren   errorx.ErrorKey = "DEPARTMENT_HAS_CHILDREN"
	ErrDepartmentHasUsers      errorx.ErrorKey = "DEPARTMENT_HAS_USERS"
	ErrInvalidParentDepartment errorx.ErrorKey = "INVALID_PARENT_DEPARTMENT"
)

var (
	ErrPostNotFound   errorx.ErrorKey = "POST_NOT_FOUND"
	ErrPostNameExists errorx.ErrorKey = "POST_NAME_EXISTS"
	ErrPostHasUsers   errorx.ErrorKey = "POST_HAS_USERS"
)

func init() {
	errorx.Register(ErrDepartmentNotFound, 404, "DEPARTMENT_NOT_FOUND", "department not found")
	errorx.Register(ErrDepartmentNameExists, 409, "DEPARTMENT_NAME_EXISTS", "department name already exists")
	errorx.Register(ErrDepartmentHasChildren, 400, "DEPARTMENT_HAS_CHILDREN", "department has children")
	errorx.Register(ErrDepartmentHasUsers, 400, "DEPARTMENT_HAS_USERS", "department has users")
	errorx.Register(ErrInvalidParentDepartment, 400, "INVALID_PARENT_DEPARTMENT", "invalid parent department")

	errorx.Register(ErrPostNotFound, 404, "POST_NOT_FOUND", "post not found")
	errorx.Register(ErrPostNameExists, 409, "POST_NAME_EXISTS", "post name already exists")
	errorx.Register(ErrPostHasUsers, 400, "POST_HAS_USERS", "post has users")
}
