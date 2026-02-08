package organization

import "time"

type Department struct {
	ID           string
	Name         string
	ParentID     string
	Sort         int32
	LeaderUserID string
	Phone        string
	Email        string
	Status       int32
	CreateBy     string
	CreateAt     time.Time
	UpdateBy     string
	UpdateAt     time.Time
	TenantID     string
	Children     []*Department
}

type Post struct {
	ID       string
	Name     string
	Code     string
	Sort     int32
	Status   int32
	Remark   string
	CreateBy string
	CreateAt time.Time
	UpdateBy string
	UpdateAt time.Time
	TenantID string
}

type ListPostsQuery struct {
	Page      int32
	PageSize  int32
	Keyword   string
	Status    *int32
	SortField string
	SortOrder string
}

type WherePostOpt struct {
	Limit     int32
	Offset    int32
	Keyword   string
	Status    *int32
	SortField string
	SortOrder string
}

type ListPostsResult struct {
	Posts      []*Post
	Total      int64
	Page       int32
	PageSize   int32
	TotalPages int32
}
