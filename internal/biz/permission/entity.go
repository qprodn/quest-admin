package permission

import "time"

type Role struct {
	ID               string
	Name             string
	Code             string
	Sort             int32
	DataScope        int32
	DataScopeDeptIDs string
	Status           int32
	Type             int32
	Remark           string
	CreateBy         string
	CreateAt         time.Time
	UpdateBy         string
	UpdateAt         time.Time
	TenantID         string
}

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
	CreateBy      string
	CreateAt      time.Time
	UpdateBy      string
	UpdateAt      time.Time
	Children      []*Menu
}

type ListRolesQuery struct {
	Page      int32
	PageSize  int32
	Keyword   string
	Status    *int32
	SortField string
	SortOrder string
}

type ListRolesResult struct {
	Roles      []*Role
	Total      int64
	Page       int32
	PageSize   int32
	TotalPages int32
}
