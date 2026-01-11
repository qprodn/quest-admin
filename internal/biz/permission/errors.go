package permission

import (
	v1 "quest-admin/api/gen/permission/v1"

	"github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrRoleNotFound      = errors.NotFound(v1.ErrorReason_ROLE_NOT_FOUND.String(), "role not found")
	ErrRoleNameExists    = errors.Conflict(v1.ErrorReason_ROLE_NAME_EXISTS.String(), "role name already exists")
	ErrRoleCodeExists    = errors.Conflict(v1.ErrorReason_ROLE_CODE_EXISTS.String(), "role code already exists")
	ErrRoleHasUsers      = errors.BadRequest(v1.ErrorReason_ROLE_HAS_USERS.String(), "role has users")
	ErrInvalidRoleStatus = errors.BadRequest(v1.ErrorReason_INVALID_ROLE_STATUS.String(), "invalid role status")
)

var (
	ErrMenuNotFound      = errors.NotFound(v1.ErrorReason_MENU_NOT_FOUND.String(), "menu not found")
	ErrMenuNameExists    = errors.Conflict(v1.ErrorReason_MENU_NAME_EXISTS.String(), "menu name already exists")
	ErrMenuHasChildren   = errors.BadRequest(v1.ErrorReason_MENU_HAS_CHILDREN.String(), "menu has children")
	ErrInvalidParentMenu = errors.BadRequest(v1.ErrorReason_INVALID_PARENT_MENU.String(), "invalid parent menu")
	ErrInvalidMenuStatus = errors.BadRequest(v1.ErrorReason_INVALID_MENU_STATUS.String(), "invalid menu status")
	ErrInvalidMenuType   = errors.BadRequest(v1.ErrorReason_INVALID_MENU_TYPE.String(), "invalid menu type")
	ErrMenuLevelExceeded = errors.BadRequest(v1.ErrorReason_MENU_LEVEL_EXCEEDED.String(), "menu level exceeded")
	ErrInvalidMenuPath   = errors.BadRequest(v1.ErrorReason_INVALID_MENU_PATH.String(), "invalid menu path")
)
