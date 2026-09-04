package http

import (
	"strconv"

	"github.com/epmp/backend/internal/modules/facility/dto"
	"github.com/epmp/backend/internal/modules/facility/service"
	mw "github.com/epmp/backend/internal/pkg/middleware"
	"github.com/epmp/backend/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

// FacilityHandler handles HTTP requests for Facility resources.
type FacilityHandler struct {
	svc *service.FacilityService
}

// NewFacilityHandler creates a new FacilityHandler.
func NewFacilityHandler(svc *service.FacilityService) *FacilityHandler {
	return &FacilityHandler{svc: svc}
}

func (h *FacilityHandler) Create(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return response.BadRequest(c, "X-Organization-ID header is required")
	}

	var req dto.CreateFacilityRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	result, err := h.svc.Create(c.Request().Context(), orgID, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, result)
}

func (h *FacilityHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	result, err := h.svc.GetByID(c.Request().Context(), id, orgID)
	if err != nil {
		return response.NotFound(c, "Facility not found")
	}

	return response.OK(c, result)
}

func (h *FacilityHandler) List(c echo.Context) error {
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

func (h *FacilityHandler) Update(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	var req dto.UpdateFacilityRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	result, err := h.svc.Update(c.Request().Context(), id, orgID, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *FacilityHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	if err := h.svc.Delete(c.Request().Context(), id, orgID); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.NoContent(c)
}
