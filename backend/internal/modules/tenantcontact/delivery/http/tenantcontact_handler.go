package http

import (
	"strconv"

	"github.com/epmp/backend/internal/modules/tenantcontact/dto"
	"github.com/epmp/backend/internal/modules/tenantcontact/service"
	"github.com/epmp/backend/internal/pkg/errs"
	mw "github.com/epmp/backend/internal/pkg/middleware"
	"github.com/epmp/backend/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

// TenantContactHandler handles HTTP requests for TenantContact resources.
type TenantContactHandler struct {
	svc *service.TenantContactService
}

// NewTenantContactHandler creates a new TenantContactHandler.
func NewTenantContactHandler(svc *service.TenantContactService) *TenantContactHandler {
	return &TenantContactHandler{svc: svc}
}

func (h *TenantContactHandler) Create(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return response.BadRequest(c, "X-Organization-ID header is required")
	}

	var req dto.CreateTenantContactRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	result, err := h.svc.Create(c.Request().Context(), orgID, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, result)
}

func (h *TenantContactHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	result, err := h.svc.GetByID(c.Request().Context(), id, orgID)
	if err != nil {
		return response.NotFound(c, "TenantContact not found")
	}

	return response.OK(c, result)
}

func (h *TenantContactHandler) List(c echo.Context) error {
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

func (h *TenantContactHandler) Update(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	var req dto.UpdateTenantContactRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	result, err := h.svc.Update(c.Request().Context(), id, orgID, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *TenantContactHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	if err := h.svc.Delete(c.Request().Context(), id, orgID); err != nil {
		if errs.IsDomainError(err, "NOT_FOUND") {
			return response.NotFound(c, "Resource not found")
		}
		return response.InternalError(c, err.Error())
	}

	return response.NoContent(c)
}
