package handler

import (
	apperror "mob/ddd-template/internal/app/error"
	"mob/ddd-template/internal/app/usecase"
	"mob/ddd-template/internal/presentation/dto"

	"github.com/gofiber/fiber/v3"
	"github.com/samber/do/v2"
)

type AuthHandler interface {
	Login(ctx fiber.Ctx) error
	Refresh(ctx fiber.Ctx) error
	Logout(ctx fiber.Ctx) error
	LogoutAll(ctx fiber.Ctx) error
}

type authHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(i do.Injector) AuthHandler {
	return &authHandler{
		authUseCase: do.MustInvoke[usecase.AuthUseCase](i),
	}
}

func (h *authHandler) Login(ctx fiber.Ctx) error {
	var req dto.AuthLoginBody
	if err := ctx.Bind().Body(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	out, err := h.authUseCase.Login(req.ToAppInput())
	if err != nil {
		if err == apperror.ErrInvalidEmailOrPassword {
			res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_DENIED_ACCESS, err.Error())
			return ctx.Status(fiber.StatusUnauthorized).JSON(res)
		}

		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(res)
	}

	ctx.Locals("userId", out.ID)
	ctx.Locals("userName", out.Name)

	res := dto.BuildResponseSuccess(dto.MESSAGE_SUCCESS_LOGIN, dto.AuthLoginOutputToResponse(out))
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *authHandler) Refresh(ctx fiber.Ctx) error {
	var req dto.RefreshTokenBody
	if err := ctx.Bind().Body(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	accessToken, err := h.authUseCase.Refresh(req.RefreshToken)
	if err != nil {
		if err == apperror.ErrRefreshTokenExpiredOrNotFound {
			res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_DENIED_ACCESS, err.Error())
			return ctx.Status(fiber.StatusUnauthorized).JSON(res)
		}

		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(res)
	}

	res := dto.BuildResponseSuccess(dto.MESSAGE_SUCCESS_REFRESH, dto.AuthRefreshResponse{
		AccessToken: accessToken,
	})
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *authHandler) Logout(ctx fiber.Ctx) error {
	var req dto.RefreshTokenBody
	if err := ctx.Bind().Body(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	if err := h.authUseCase.Logout(req.RefreshToken); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(res)
	}

	res := dto.BuildResponseSuccess(dto.MESSAGE_SUCCESS_LOGOUT, nil)
	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (h *authHandler) LogoutAll(ctx fiber.Ctx) error {
	var req dto.UserIDURI
	if err := ctx.Bind().URI(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(res)
	}

	if err := h.authUseCase.LogoutAll(req.ID); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROCESS_REQUEST, err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(res)
	}

	res := dto.BuildResponseSuccess(dto.MESSAGE_SUCCESS_LOGOUT_ALL, nil)
	return ctx.Status(fiber.StatusOK).JSON(res)
}
