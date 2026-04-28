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

type UserIDURI struct {
	ID string `uri:"id" binding:"required"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`

	CreatedAt string `json:"created_at"`
}

func UserOutputToResponse(out *appDto.UserOutput) *UserResponse {
	return &UserResponse{
		ID:        out.ID,
		Name:      out.Name,
		Email:     out.Email,
		CreatedAt: out.CreatedAt,
	}
}

type PaginatedUserResponse struct {
	Data []*UserResponse `json:"data"`
	Meta *Meta           `json:"meta"`
}

func PaginatedUserOutputToResponse(out *appDto.PaginatedOutput[*appDto.UserOutput]) *PaginatedUserResponse {
	data := make([]*UserResponse, 0, len(out.Data))
	for _, user := range out.Data {
		data = append(data, UserOutputToResponse(user))
	}

	return &PaginatedUserResponse{
		Data: data,
		Meta: &Meta{
			Limit:     out.Limit,
			Page:      out.Page,
			TotalData: out.TotalData,
			TotalPage: out.TotalPage,
		},
	}
}
