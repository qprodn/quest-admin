package tenant

import "time"

type Tenant struct {
	ID            string
	Name          string
	ContactUserID string
	ContactName   string
	ContactMobile string
	Status        int32
	Website       string
	PackageID     string
	ExpireTime    time.Time
	AccountCount  int32
	CreateBy      string
	CreateAt      time.Time
	UpdateBy      string
	UpdateAt      time.Time
}

type TenantPackage struct {
	ID       string
	Name     string
	Status   int32
	Remark   string
	MenuIDs  string
	CreateBy string
	CreateAt time.Time
	UpdateBy string
	UpdateAt time.Time
}

type ListTenantsQuery struct {
	Page      int32
	PageSize  int32
	Keyword   string
	Status    *int32
	SortField string
	SortOrder string
}

type ListTenantsResult struct {
	Tenants    []*Tenant
	Total      int64
	Page       int32
	PageSize   int32
	TotalPages int32
}

type ListPackagesQuery struct {
	Page      int32
	PageSize  int32
	Keyword   string
	Status    *int32
	SortField string
	SortOrder string
}

type ListPackagesResult struct {
	Packages   []*TenantPackage
	Total      int64
	Page       int32
	PageSize   int32
	TotalPages int32
}
