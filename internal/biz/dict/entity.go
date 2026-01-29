package dict

import "time"

type DictType struct {
	ID       string
	Name     string
	Code     string
	Status   int32
	Sort     int32
	Remark   string
	CreateBy string
	CreateAt time.Time
	UpdateBy string
	UpdateAt time.Time
	TenantID string
}

type ListDictTypesQuery struct {
	Page      int32
	PageSize  int32
	Keyword   string
	Status    *int32
	SortField string
	SortOrder string
}

type ListDictTypesResult struct {
	DictTypes  []*DictType
	Total      int64
	Page       int32
	PageSize   int32
	TotalPages int32
}

type DictData struct {
	ID         string
	DictTypeID string
	Label      string
	Value      string
	Sort       int32
	Status     int32
	CSSClass   string
	IsDefault  bool
	Remark     string
	CreateBy   string
	CreateAt   time.Time
	UpdateBy   string
	UpdateAt   time.Time
	TenantID   string
}

type ListDictDataQuery struct {
	Page       int32
	PageSize   int32
	DictTypeID string
	Keyword    string
	Status     *int32
	SortField  string
	SortOrder  string
}

type ListDictDataResult struct {
	DictData   []*DictData
	Total      int64
	Page       int32
	PageSize   int32
	TotalPages int32
}
