package handler

import (
	apperror "mob/ddd-template/internal/app/error"
	"mob/ddd-template/internal/app/usecase"
	"mob/ddd-template/internal/presentation/dto"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/samber/do/v2"
)

type UserHandler interface {
	Create(ctx fiber.Ctx) error
	GetAll(ctx fiber.Ctx) error
	GetById(ctx fiber.Ctx) error
	Delete(ctx fiber.Ctx) error
}

type userHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(i do.Injector) UserHandler {
	return &userHandler{
		userUseCase: do.MustInvoke[usecase.UserUseCase](i),
	}
}

func (h *userHandler) Create(ctx fiber.Ctx) error {
	var req dto.CreateUserBody
	if err := ctx.Bind().Body(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	out, err := h.userUseCase.Create(req.ToAppInput())
	if err != nil {
		if err == apperror.ErrEmailAlreadyInUse {
			res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
			return ctx.Status(fiber.StatusConflict).JSON(res)
		}

		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(res)
	}

	res := dto.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_DATA, dto.UserOutputToResponse(out))
	return ctx.Status(fiber.StatusCreated).JSON(res)
}

func (h *userHandler) GetAll(ctx fiber.Ctx) error {
	var req dto.PaginateQuery
	if err := ctx.Bind().Query(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	out, err := h.userUseCase.GetAll(req.ToAppInput())
	if err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(http.StatusInternalServerError).JSON(res)
	}

	dtoRes := dto.PaginatedUserOutputToResponse(out)
	res := dto.BuildPaginatedResponseSuccess(dto.MESSAGE_SUCCESS_GET_DATA, dtoRes.Data, dtoRes.Meta)
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *userHandler) GetById(ctx fiber.Ctx) error {
	var req dto.UserIDURI
	if err := ctx.Bind().URI(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	out, err := h.userUseCase.GetById(req.ID)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err == apperror.ErrUserNotFound {
			status = fiber.StatusNotFound
		}

		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(status).JSON(res)
	}

	res := dto.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DATA, dto.UserOutputToResponse(out))
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *userHandler) Delete(ctx fiber.Ctx) error {
	var req dto.UserIDURI
	if err := ctx.Bind().URI(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	if err := h.userUseCase.Delete(req.ID); err != nil {
		status := fiber.StatusInternalServerError
		if err == apperror.ErrUserNotFound {
			status = fiber.StatusNotFound
		}

		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(status).JSON(res)
	}

	res := dto.BuildResponseSuccess(dto.MESSAGE_SUCCESS_DELETE_DATA, nil)
	return ctx.Status(fiber.StatusNoContent).JSON(res)
}
