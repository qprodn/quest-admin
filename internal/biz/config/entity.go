package config

import "time"

type Config struct {
	ID        string
	Name      string
	Key       string
	Value     string
	Status    int32
	CreateBy  string
	CreateAt  time.Time
	UpdateBy  string
	UpdateAt  time.Time
	TenantID  string
}

type ListConfigsQuery struct {
	Page      int32
	PageSize  int32
	Name      string
	Key       string
	Status    *int32
	SortField string
	SortOrder string
}

type WhereConfigOpt struct {
	Limit     int32
	Offset    int32
	Name      string
	Key       string
	Status    *int32
	SortField string
	SortOrder string
}

type ListConfigsResult struct {
	Configs    []*Config
	Total      int64
	Page       int32
	PageSize   int32
	TotalPages int32
}

type DeleteConfigBO struct {
	ConfigID string
}

type UpdateStatusBO struct {
	ConfigID string
	Status   int32
}
