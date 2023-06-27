package pagination

type Pagination struct {
	Page     int8 `form:"page"`
	PageSize int8 `form:"page_size"`
}
