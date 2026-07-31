package http

import (
	"strconv"

	"github.com/epmp/backend/internal/modules/iam/dto"
	"github.com/epmp/backend/internal/modules/iam/service"
	"github.com/epmp/backend/internal/pkg/errs"
	"github.com/epmp/backend/internal/pkg/response"
	"github.com/labstack/echo/v4"
)

// UserHandler handles user management HTTP endpoints.
type UserHandler struct {
	userSvc *service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func (h *UserHandler) Create(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	res, err := h.userSvc.Create(c.Request().Context(), &req)
	if err != nil {
		if errs.IsDomainError(err, "ALREADY_EXISTS") {
			return response.Conflict(c, "email already registered")
		}
		return response.InternalError(c, err.Error())
	}
	return response.Created(c, res)
}

func (h *UserHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	res, err := h.userSvc.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.NotFound(c, "user not found")
	}
	return response.OK(c, res)
}

func (h *UserHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage == 0 {
		perPage = 20
	}

	res, err := h.userSvc.List(c.Request().Context(), page, perPage)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.OK(c, res)
}

func (h *UserHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	res, err := h.userSvc.Update(c.Request().Context(), id, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.OK(c, res)
}

func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.userSvc.Delete(c.Request().Context(), id); err != nil {
		return response.NotFound(c, "user not found")
	}
	return response.NoContent(c)
}

func (h *UserHandler) AssignRoles(c echo.Context) error {
	id := c.Param("id")
	var req dto.AssignRolesRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	res, err := h.userSvc.AssignRoles(c.Request().Context(), id, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.OK(c, res)
}
