package dto

import "mob/ddd-template/internal/app/dto"

type PaginateQuery struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Sort     string `form:"sort"`
	SortBy   string `form:"sortBy"`
	Filter   string `form:"filter"`
	FilterBy string `form:"filterBy"`
}

func (q *PaginateQuery) ToAppInput() *dto.PaginateInput {
	return &dto.PaginateInput{
		Page:     q.Page,
		Limit:    q.Limit,
		Sort:     q.Sort,
		SortBy:   q.SortBy,
		Filter:   q.Filter,
		FilterBy: q.FilterBy,
	}
}
