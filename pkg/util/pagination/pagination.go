package pagination

func GetOffset(page, pageSize int32) int32 {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return (page - 1) * pageSize
}

func GetTotalPages(total, pageSize int64) int32 {
	if total < 0 {
		total = 0
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return int32((total + pageSize - 1) / pageSize)
}
