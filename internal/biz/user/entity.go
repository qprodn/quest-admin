package user

import "time"

type User struct {
	ID        string
	Username  string
	Password  string
	Nickname  string
	Email     string
	Mobile    string
	Sex       int32
	Avatar    string
	Status    int32
	Remark    string
	LoginIP   string
	LoginDate time.Time
	CreateBy  string
	CreateAt  time.Time
	UpdateBy  string
	UpdateAt  time.Time
	TenantID  string
}

type UpdatePasswordBO struct {
	UserID      string
	OldPassword string
	NewPassword string
}

type UpdateStatusBO struct {
	UserID string
	Status int32
}

type UpdateLoginInfoBO struct {
	UserID    string
	LoginIP   string
	LoginDate time.Time
}

type ManageUserPostsBO struct {
	UserID    string
	PostIDs   []string
	Operation string
}

type ListUsersQuery struct {
	Page      int32
	PageSize  int32
	Keyword   string
	Status    *int32
	Sex       *int32
	SortField string
	SortOrder string
}

type ListUsersResult struct {
	Users      []*User
	Total      int64
	Page       int32
	PageSize   int32
	TotalPages int32
}

type ManageUserRolesBO struct {
	UserID    string
	RoleIDs   []string
	Operation string
}
