package http

import (
	"github.com/epmp/backend/internal/modules/iam/application/dto"
	"github.com/epmp/backend/internal/modules/iam/application/service"
	"github.com/epmp/backend/internal/shared"
	"github.com/epmp/backend/internal/shared/response"
	"github.com/labstack/echo/v4"
)

// RoleHandler handles role & permission HTTP endpoints.
type RoleHandler struct {
	roleSvc *service.RoleService
}

// NewRoleHandler creates a new RoleHandler.
func NewRoleHandler(roleSvc *service.RoleService) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc}
}

func (h *RoleHandler) Create(c echo.Context) error {
	var req dto.CreateRoleRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	res, err := h.roleSvc.Create(c.Request().Context(), &req)
	if err != nil {
		if shared.IsDomainError(err, "ALREADY_EXISTS") {
			return response.Conflict(c, "role name already exists")
		}
		return response.InternalError(c, err.Error())
	}
	return response.Created(c, res)
}

func (h *RoleHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	res, err := h.roleSvc.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.NotFound(c, "role not found")
	}
	return response.OK(c, res)
}

func (h *RoleHandler) List(c echo.Context) error {
	res, err := h.roleSvc.List(c.Request().Context())
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.OK(c, res)
}

func (h *RoleHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req dto.UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	res, err := h.roleSvc.Update(c.Request().Context(), id, &req)
	if err != nil {
		if shared.IsDomainError(err, "SYSTEM_ROLE") {
			return response.Forbidden(c, "system roles cannot be modified")
		}
		return response.InternalError(c, err.Error())
	}
	return response.OK(c, res)
}

func (h *RoleHandler) SetPermissions(c echo.Context) error {
	id := c.Param("id")
	var req dto.SetPermissionsRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	res, err := h.roleSvc.SetPermissions(c.Request().Context(), id, &req)
	if err != nil {
		if shared.IsDomainError(err, "SYSTEM_ROLE") {
			return response.Forbidden(c, "system roles cannot be modified")
		}
		return response.InternalError(c, err.Error())
	}
	return response.OK(c, res)
}

func (h *RoleHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.roleSvc.Delete(c.Request().Context(), id); err != nil {
		if shared.IsDomainError(err, "SYSTEM_ROLE") {
			return response.Forbidden(c, "system roles cannot be deleted")
		}
		return response.NotFound(c, "role not found")
	}
	return response.NoContent(c)
}

func (h *RoleHandler) ListPermissions(c echo.Context) error {
	res, err := h.roleSvc.ListPermissions(c.Request().Context())
	if err != nil {
		return response.InternalError(c, err.Error())
	}
	return response.OK(c, res)
}
