package permission

import (
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrRoleNotFound      = errors.NotFound("ROLE_NOT_FOUND", "role not found")
	ErrRoleNameExists    = errors.Conflict("ROLE_NAME_EXISTS", "role name already exists")
	ErrRoleCodeExists    = errors.Conflict("ROLE_CODE_EXISTS", "role code already exists")
	ErrRoleHasUsers      = errors.BadRequest("ROLE_HAS_USERS", "role has users")
	ErrInvalidRoleStatus = errors.BadRequest("INVALID_ROLE_STATUS", "invalid role status")
)

var (
	ErrMenuNotFound      = errors.NotFound("MENU_NOT_FOUND", "menu not found")
	ErrMenuNameExists    = errors.Conflict("MENU_NAME_EXISTS", "menu name already exists")
	ErrMenuHasChildren   = errors.BadRequest("MENU_HAS_CHILDREN", "menu has children")
	ErrInvalidParentMenu = errors.BadRequest("INVALID_PARENT_MENU", "invalid parent menu")
	ErrInvalidMenuStatus = errors.BadRequest("INVALID_MENU_STATUS", "invalid menu status")
	ErrInvalidMenuType   = errors.BadRequest("INVALID_MENU_TYPE", "invalid menu type")
	ErrMenuLevelExceeded = errors.BadRequest("MENU_LEVEL_EXCEEDED", "menu level exceeded")
	ErrInvalidMenuPath   = errors.BadRequest("INVALID_MENU_PATH", "invalid menu path")
)

var (
	ErrInternalServer = errors.InternalServer("UNKNOWN_SYSTEM_ERROR", "UNKNOWN_SYSTEM_ERROR")
)
