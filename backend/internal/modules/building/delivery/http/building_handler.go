package http

import (
	"strconv"

	"github.com/epmp/backend/internal/modules/building/dto"
	"github.com/epmp/backend/internal/modules/building/service"
	mw "github.com/epmp/backend/internal/pkg/middleware"
	"github.com/epmp/backend/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

// BuildingHandler handles HTTP requests for Building resources.
type BuildingHandler struct {
	svc *service.BuildingService
}

// NewBuildingHandler creates a new BuildingHandler.
func NewBuildingHandler(svc *service.BuildingService) *BuildingHandler {
	return &BuildingHandler{svc: svc}
}

func (h *BuildingHandler) Create(c echo.Context) error {
	var req dto.CreateBuildingRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Name == "" {
		return response.BadRequest(c, "name is required")
	}
	if req.PropertyId == "" {
		return response.BadRequest(c, "property_id is required")
	}

	orgID := mw.GetOrgID(c)

	result, err := h.svc.Create(c.Request().Context(), &req, orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, result)
}

func (h *BuildingHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	result, err := h.svc.GetByID(c.Request().Context(), id, orgID)
	if err != nil {
		return response.NotFound(c, "Building not found")
	}

	return response.OK(c, result)
}

func (h *BuildingHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage == 0 {
		perPage = 20
	}
	search := c.QueryParam("search")
	propertyId := c.QueryParam("property_id")
	orgID := mw.GetOrgID(c)

	result, err := h.svc.List(c.Request().Context(), page, perPage, search, propertyId, orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *BuildingHandler) Update(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	var req dto.UpdateBuildingRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	result, err := h.svc.Update(c.Request().Context(), id, orgID, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *BuildingHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.NoContent(c)
}
