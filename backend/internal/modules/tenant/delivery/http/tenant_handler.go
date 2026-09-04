package http

import (
	"strconv"

	"github.com/epmp/backend/internal/modules/tenant/dto"
	"github.com/epmp/backend/internal/modules/tenant/service"
	mw "github.com/epmp/backend/internal/pkg/middleware"
	"github.com/epmp/backend/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

// TenantHandler handles HTTP requests for Tenant resources.
type TenantHandler struct {
	svc *service.TenantService
}

// NewTenantHandler creates a new TenantHandler.
func NewTenantHandler(svc *service.TenantService) *TenantHandler {
	return &TenantHandler{svc: svc}
}

func (h *TenantHandler) Create(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return response.BadRequest(c, "X-Organization-ID header is required")
	}

	var req dto.CreateTenantRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.FullName == "" {
		return response.BadRequest(c, "full_name is required")
	}
	if req.Email == "" {
		return response.BadRequest(c, "email is required")
	}

	result, err := h.svc.Create(c.Request().Context(), orgID, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, result)
}

func (h *TenantHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	result, err := h.svc.GetByID(c.Request().Context(), id, orgID)
	if err != nil {
		return response.NotFound(c, "Tenant not found")
	}

	return response.OK(c, result)
}

func (h *TenantHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage == 0 {
		perPage = 20
	}
	search := c.QueryParam("search")
	orgID := mw.GetOrgID(c)

	result, err := h.svc.List(c.Request().Context(), page, perPage, search, orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *TenantHandler) Update(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	var req dto.UpdateTenantRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.FullName == "" {
		return response.BadRequest(c, "full_name is required")
	}
	if req.Email == "" {
		return response.BadRequest(c, "email is required")
	}

	result, err := h.svc.Update(c.Request().Context(), id, orgID, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *TenantHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	if err := h.svc.Delete(c.Request().Context(), id, orgID); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.NoContent(c)
}
