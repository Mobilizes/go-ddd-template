package handler

import (
	"mob/ddd-template/internal/app"
	"mob/ddd-template/internal/presentation/dto"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/samber/do/v2"
)

type UserHandler interface {
	Create(c fiber.Ctx) error
	GetAll(c fiber.Ctx) error
	GetById(c fiber.Ctx) error
	Delete(c fiber.Ctx) error
}

type userPresentation struct {
	userUseCase app.UserUseCase
}

func NewUserPresentation(i do.Injector) UserHandler {
	return &userPresentation{
		userUseCase: do.MustInvoke[app.UserUseCase](i),
	}
}

func RegisterUserRoutes(r *fiber.App, p UserHandler) {
	userGroup := r.Group("/api/user")
	{
		userGroup.Post("", p.Create)
		userGroup.Get("", p.GetAll)
		userGroup.Get("/:id", p.GetById)
		userGroup.Delete("/:id", p.Delete)
	}
}

func (p *userPresentation) Create(c fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	out, err := p.userUseCase.Create(req.ToAppInput())
	if err != nil {
		if err == app.ErrEmailAlreadyInUse {
			res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
			return c.Status(fiber.StatusConflict).JSON(res)
		}

		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}

	res := dto.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_DATA, dto.UserOutputToResponse(out))
	return c.Status(fiber.StatusCreated).JSON(res)
}

func (p *userPresentation) GetAll(c fiber.Ctx) error {
	var req dto.PaginateQuery
	if err := c.Bind().Query(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	out, err := p.userUseCase.GetAll(req.ToAppInput())
	if err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(res)
	}

	dtoRes := dto.PaginatedUserOutputToResponse(out)
	res := dto.BuildPaginatedResponseSuccess(dto.MESSAGE_SUCCESS_GET_DATA, dtoRes.Data, dtoRes.Meta)
	return c.Status(fiber.StatusOK).JSON(res)
}

func (p *userPresentation) GetById(c fiber.Ctx) error {
	var req dto.UserIDURI
	if err := c.Bind().URI(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	out, err := p.userUseCase.GetById(req.ID)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err == app.ErrUserNotFound {
			status = fiber.StatusNotFound
		}

		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		return c.Status(status).JSON(res)
	}

	res := dto.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DATA, dto.UserOutputToResponse(out))
	return c.Status(fiber.StatusOK).JSON(res)
}

func (p *userPresentation) Delete(c fiber.Ctx) error {
	var req dto.UserIDURI
	if err := c.Bind().URI(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	if err := p.userUseCase.Delete(req.ID); err != nil {
		status := fiber.StatusInternalServerError
		if err == app.ErrUserNotFound {
			status = fiber.StatusNotFound
		}

		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		return c.Status(status).JSON(res)
	}

	res := dto.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_DATA, nil)
	return c.Status(fiber.StatusNoContent).JSON(res)
}
