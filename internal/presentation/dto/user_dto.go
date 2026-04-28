package dto

import appDto "mob/ddd-template/internal/app/dto"

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (r *CreateUserRequest) ToAppInput() *appDto.UserCreateInput {
	return &appDto.UserCreateInput{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}

type GetUsersQuery struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Sort     string `form:"sort"`
	SortBy   string `form:"sortBy"`
	Filter   string `form:"filter"`
	FilterBy string `form:"filterBy"`
}

func (q *GetUsersQuery) ToAppInput() *appDto.PaginateInput {
	return &appDto.PaginateInput{
		Page:     q.Page,
		Limit:    q.Limit,
		Sort:     q.Sort,
		SortBy:   q.SortBy,
		Filter:   q.Filter,
		FilterBy: q.FilterBy,
	}
}

type UserIDURI struct {
	ID string `uri:"id" binding:"required"`
}

type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func UserOutputToResponse(out *appDto.UserOutput) *UserResponse {
	return &UserResponse{
		Name:  out.Name,
		Email: out.Email,
	}
}

type PaginatedUserResponse struct {
	Data      []*UserResponse `json:"data"`
	Limit     int             `json:"limit"`
	Page      int             `json:"page"`
	TotalData int64           `json:"totalData"`
	TotalPage int             `json:"totalPage"`
}

func PaginatedUserOutputToResponse(out *appDto.PaginatedOutput[*appDto.UserOutput]) *PaginatedUserResponse {
	data := make([]*UserResponse, 0, len(out.Data))
	for _, user := range out.Data {
		data = append(data, UserOutputToResponse(user))
	}

	return &PaginatedUserResponse{
		Data:      data,
		Limit:     out.Limit,
		Page:      out.Page,
		TotalData: out.TotalData,
		TotalPage: out.TotalPage,
	}
}
