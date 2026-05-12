package usecase

import (
	"mob/ddd-template/internal/app/dto"
	apperror "mob/ddd-template/internal/app/error"
	"mob/ddd-template/internal/app/port"
	"mob/ddd-template/internal/domain/entity"
	"mob/ddd-template/internal/domain/repository"
	vo "mob/ddd-template/internal/domain/valueobject"
	"slices"

	"github.com/google/uuid"
	"github.com/samber/do/v2"
)

type UserUseCase interface {
	Create(req *dto.UserCreateInput) (*dto.UserOutput, error)
	GetAll(req *dto.PaginateInput) (*dto.PaginatedOutput[*dto.UserOutput], error)
	GetById(id string) (*dto.UserOutput, error)
	Delete(id string) error
}

type userUseCase struct {
	userRepository repository.UserRepository
	hasher         port.Hasher
	tokenGenerator port.TokenGenerator
	unitOfWork     port.UnitOfWork
}

func NewUserUseCase(i do.Injector) UserUseCase {
	return &userUseCase{
		userRepository: do.MustInvoke[repository.UserRepository](i),
		hasher:         do.MustInvoke[port.Hasher](i),
		tokenGenerator: do.MustInvoke[port.TokenGenerator](i),
		unitOfWork:     do.MustInvoke[port.UnitOfWork](i),
	}
}

func (uc *userUseCase) Create(req *dto.UserCreateInput) (*dto.UserOutput, error) {
	_, err := uc.userRepository.GetByEmail(req.Email)
	if err == nil {
		return &dto.UserOutput{}, apperror.ErrEmailAlreadyInUse
	}

	hashedPassword, err := uc.hasher.RandomHash(req.Password)
	if err != nil {
		return &dto.UserOutput{}, err
	}

	user := entity.NewUser(uuid.NewString(), req.Name, req.Email, hashedPassword)

	if err := uc.userRepository.Create(user); err != nil {
		return &dto.UserOutput{}, err
	}

	return dto.UserEntityToOutput(user), nil
}

func (uc *userUseCase) GetAll(req *dto.PaginateInput) (*dto.PaginatedOutput[*dto.UserOutput], error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	sort := req.Sort
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	var allowedSort = []string{"name", "email", "created_at", "updated_at"}
	sortBy := req.SortBy
	if !slices.Contains(allowedSort, sortBy) {
		sortBy = "id"
	}

	var allowedFilter = []string{"name", "email", "created_at", "updated_at"}
	filterBy := req.FilterBy
	if !slices.Contains(allowedFilter, filterBy) {
		filterBy = ""
	}

	opts := vo.PaginateOptions{
		Page:     max(0, req.Page),
		Limit:    limit,
		Sort:     sort,
		SortBy:   sortBy,
		Filter:   req.Filter,
		FilterBy: filterBy,
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
	user, err := uc.userRepository.GetById(id)
	if err != nil {
		return &dto.UserOutput{}, apperror.ErrUserNotFound
	}

	return dto.UserEntityToOutput(user), nil
}

func (uc *userUseCase) Delete(id string) error {
	return uc.unitOfWork.Transaction(func(repos port.UnitOfWorkRepositories) error {
		userRepo := repos.Users()

		_, err := userRepo.GetById(id)
		if err != nil {
			return apperror.ErrUserNotFound
		}

		err = userRepo.Delete(id)
		if err != nil {
			return err
		}

		return nil
	})
}
