package auth

import "time"

// LoginRequest 登录请求
type LoginRequest struct {
	LoginType int32
	Username  string
	Password  string
	Device    string
}

// LoginReply 登录响应
type LoginReply struct {
	Token string
	User  *User
}

type GenerateTokenBO struct {
	UserID string
	Device string
}

// User 用户信息
type User struct {
	ID       string
	Username string
	Nickname string
	Email    string
	Mobile   string
	Sex      int32
	Avatar   string
	Status   int32
	Remark   string
	CreateAt time.Time
}

// GetPermissionInfoReply 获取权限信息响应
type GetPermissionInfoReply struct {
	User        *User
	Roles       []string
	Permissions []string
	Menus       []*Menu
}

// Menu 菜单信息
type Menu struct {
	ID            string
	Name          string
	Permission    string
	Type          int32
	Sort          int32
	ParentID      string
	Path          string
	Icon          string
	Component     string
	ComponentName string
	Status        int32
	Visible       bool
	KeepAlive     bool
	AlwaysShow    bool
	Children      []*Menu
}
