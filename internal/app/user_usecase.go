package app

import (
	"errors"
	"mob/ddd-template/internal/app/dto"
	"mob/ddd-template/internal/domain/entity"
	"mob/ddd-template/internal/domain/repository"
	"mob/ddd-template/internal/domain/service"
	vo "mob/ddd-template/internal/domain/valueobject"

	"github.com/samber/do/v2"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyInUse  = errors.New("email already in use")
)

type UserUseCase interface {
	Create(req *dto.UserCreateInput) (*dto.UserOutput, error)
	GetAll(req *dto.PaginateInput) (*dto.PaginatedOutput[*dto.UserOutput], error)
	GetById(id string) (*dto.UserOutput, error)
	Delete(id string) error
}

type userUseCase struct {
	userRepository repository.UserRepository
	hasher         service.PasswordHasher
	tokenGenerator service.TokenGenerator
}

func NewUserUseCase(i do.Injector) UserUseCase {
	return &userUseCase{
		userRepository: do.MustInvoke[repository.UserRepository](i),
		hasher:         do.MustInvoke[service.PasswordHasher](i),
		tokenGenerator: do.MustInvoke[service.TokenGenerator](i),
	}
}

func (uc *userUseCase) Create(req *dto.UserCreateInput) (*dto.UserOutput, error) {
	_, err := uc.userRepository.GetByEmail(req.Email)
	if err == nil {
		return &dto.UserOutput{}, ErrEmailAlreadyInUse
	}

	hashedPassword, err := uc.hasher.HashPassword(req.Password)
	if err != nil {
		return &dto.UserOutput{}, err
	}

	user := entity.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hashedPassword,
	}

	if err := uc.userRepository.Create(&user); err != nil {
		return &dto.UserOutput{}, err
	}

	return dto.UserEntityToOutput(&user), nil
}

func (uc *userUseCase) GetAll(req *dto.PaginateInput) (*dto.PaginatedOutput[*dto.UserOutput], error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	sort := req.Sort
	if sort == "" {
		sort = "asc"
	}

	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "id"
	}

	opts := vo.PaginateOptions{
		Page:     max(0, req.Page),
		Limit:    limit,
		Sort:     sort,
		SortBy:   req.SortBy,
		Filter:   req.Filter,
		FilterBy: req.FilterBy,
	}

	result, err := uc.userRepository.GetAll(&opts)
	if err != nil {
		return nil, err
	}

	var userOutputs []*dto.UserOutput
	for _, user := range result.Data {
		userOutputs = append(userOutputs, dto.UserEntityToOutput(user))
	}

	return &dto.PaginatedOutput[*dto.UserOutput]{
		Data:      userOutputs,
		Page:      opts.Page,
		Limit:     opts.Limit,
		TotalData: result.TotalData,
		TotalPage: result.TotalPage,
	}, nil
}

func (uc *userUseCase) GetById(id string) (*dto.UserOutput, error) {
	user, err := uc.userRepository.GetByID(id)
	if err != nil {
		return &dto.UserOutput{}, ErrUserNotFound
	}

	return dto.UserEntityToOutput(user), nil
}

func (uc *userUseCase) Delete(id string) error {
	return uc.userRepository.Delete(id)
}
