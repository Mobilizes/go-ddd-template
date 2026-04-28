package presentation

import (
	"net/http"

	"mob/ddd-template/internal/app"
	"mob/ddd-template/internal/presentation/dto"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type UserPresentation interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	GetById(c *gin.Context)
	Delete(c *gin.Context)
}

type userPresentation struct {
	userUseCase app.UserUseCase
}

func NewUserPresentation(i do.Injector) UserPresentation {
	return &userPresentation{
		userUseCase: do.MustInvoke[app.UserUseCase](i),
	}
}

func RegisterUserRoutes(r *gin.Engine, p UserPresentation) {
	userGroup := r.Group("/api/users")
	{
		userGroup.POST("", p.Create)
		userGroup.GET("", p.GetAll)
		userGroup.GET("/:id", p.GetById)
		userGroup.DELETE("/:id", p.Delete)
	}
}

func (p *userPresentation) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	out, err := p.userUseCase.Create(req.ToAppInput())
	if err != nil {
		if err == app.ErrEmailAlreadyInUse {
			res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
			c.JSON(http.StatusConflict, res)
			return
		}

		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := dto.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_DATA, dto.UserOutputToResponse(out))
	c.JSON(http.StatusCreated, res)
}

func (p *userPresentation) GetAll(c *gin.Context) {
	var req dto.GetUsersQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	out, err := p.userUseCase.GetAll(req.ToAppInput())
	if err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := dto.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DATA, dto.PaginatedUserOutputToResponse(out))
	c.JSON(http.StatusOK, res)
}

func (p *userPresentation) GetById(c *gin.Context) {
	var req dto.UserIDURI
	if err := c.ShouldBindUri(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	out, err := p.userUseCase.GetById(req.ID)
	if err != nil {
		if err == app.ErrUserNotFound {
			res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
			c.JSON(http.StatusNotFound, res)
			return
		}
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := dto.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_DATA, dto.UserOutputToResponse(out))
	c.JSON(http.StatusOK, res)
}

func (p *userPresentation) Delete(c *gin.Context) {
	var req dto.UserIDURI
	if err := c.ShouldBindUri(&req); err != nil {
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := p.userUseCase.Delete(req.ID); err != nil {
		if err == app.ErrUserNotFound {
			res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
			c.JSON(http.StatusNotFound, res)
			return
		}
		res := dto.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error())
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.JSON(http.StatusNoContent, dto.MESSAGE_SUCCESS_DELETE_DATA)
}
