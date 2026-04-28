package vo

type PaginateOptions struct {
	Page     int
	Limit    int
	Sort     string
	SortBy   string
	Filter   string
	FilterBy string
}

type PaginatedResult[T any] struct {
	Data      []T
	Limit     int
	Page      int
	TotalData int64
	TotalPage int
}
