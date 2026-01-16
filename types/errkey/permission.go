package errkey

import "quest-admin/pkg/errorx"

var (
	ErrMenuNotFound      errorx.ErrorKey = "MENU_NOT_FOUND"
	ErrMenuNameExists    errorx.ErrorKey = "MENU_NAME_EXISTS"
	ErrMenuHasChildren   errorx.ErrorKey = "MENU_HAS_CHILDREN"
	ErrInvalidParentMenu errorx.ErrorKey = "INVALID_PARENT_MENU"
	ErrInvalidMenuStatus errorx.ErrorKey = "INVALID_MENU_STATUS"
	ErrInvalidMenuType   errorx.ErrorKey = "INVALID_MENU_TYPE"
	ErrMenuLevelExceeded errorx.ErrorKey = "MENU_LEVEL_EXCEEDED"
	ErrInvalidMenuPath   errorx.ErrorKey = "INVALID_MENU_PATH"
)

var (
	ErrRoleNotFound      errorx.ErrorKey = "ROLE_NOT_FOUND"
	ErrRoleNameExists    errorx.ErrorKey = "ROLE_NAME_EXISTS"
	ErrRoleCodeExists    errorx.ErrorKey = "ROLE_CODE_EXISTS"
	ErrRoleHasUsers      errorx.ErrorKey = "ROLE_HAS_USERS"
	ErrInvalidRoleStatus errorx.ErrorKey = "INVALID_ROLE_STATUS"
)

func init() {
	errorx.Register(ErrRoleNotFound, 404, "ROLE_NOT_FOUND", "role not found")
	errorx.Register(ErrRoleNameExists, 409, "ROLE_NAME_EXISTS", "role name already exists")
	errorx.Register(ErrRoleCodeExists, 409, "ROLE_CODE_EXISTS", "role code already exists")
	errorx.Register(ErrRoleHasUsers, 400, "ROLE_HAS_USERS", "role has users")
	errorx.Register(ErrInvalidRoleStatus, 400, "INVALID_ROLE_STATUS", "invalid role status")

	errorx.Register(ErrMenuNotFound, 404, "MENU_NOT_FOUND", "menu not found")
	errorx.Register(ErrMenuNameExists, 409, "MENU_NAME_EXISTS", "menu name already exists")
	errorx.Register(ErrMenuHasChildren, 400, "MENU_HAS_CHILDREN", "menu has children")
	errorx.Register(ErrInvalidParentMenu, 400, "INVALID_PARENT_MENU", "invalid parent menu")
	errorx.Register(ErrInvalidMenuStatus, 400, "INVALID_MENU_STATUS", "invalid menu status")
	errorx.Register(ErrInvalidMenuType, 400, "INVALID_MENU_TYPE", "invalid menu type")
	errorx.Register(ErrMenuLevelExceeded, 400, "MENU_LEVEL_EXCEEDED", "menu level exceeded")
	errorx.Register(ErrInvalidMenuPath, 400, "INVALID_MENU_PATH", "invalid menu path")
}
