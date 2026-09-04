package http

import (
	"strconv"

	"github.com/epmp/backend/internal/modules/organization/dto"
	"github.com/epmp/backend/internal/modules/organization/service"
	"github.com/epmp/backend/internal/pkg/response"
	mw "github.com/epmp/backend/internal/pkg/middleware"

	"github.com/labstack/echo/v4"
)

// OrganizationHandler handles HTTP requests for Organization resources.
type OrganizationHandler struct {
	svc *service.OrganizationService
}

// NewOrganizationHandler creates a new OrganizationHandler.
func NewOrganizationHandler(svc *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{svc: svc}
}

func (h *OrganizationHandler) Create(c echo.Context) error {
	var req dto.CreateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.Name == "" {
		return response.BadRequest(c, "name is required")
	}

	// Extract creator's user ID from JWT context
	createdBy := mw.GetUserID(c)
	if createdBy == "" {
		return response.Unauthorized(c, "unauthorized")
	}

	result, err := h.svc.Create(c.Request().Context(), &req, createdBy)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, result)
}

// ListMine returns organizations where the authenticated user is a member.
func (h *OrganizationHandler) ListMine(c echo.Context) error {
	userID := mw.GetUserID(c)
	if userID == "" {
		return response.Unauthorized(c, "unauthorized")
	}

	result, err := h.svc.ListByUser(c.Request().Context(), userID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *OrganizationHandler) GetByID(c echo.Context) error {
	id := c.Param("id")

	result, err := h.svc.GetByID(c.Request().Context(), id)
	if err != nil {
		return response.NotFound(c, "Organization not found")
	}

	return response.OK(c, result)
}

func (h *OrganizationHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage == 0 {
		perPage = 20
	}
	search := c.QueryParam("search")

	result, err := h.svc.List(c.Request().Context(), page, perPage, search)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *OrganizationHandler) Update(c echo.Context) error {
	id := c.Param("id")

	var req dto.UpdateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	result, err := h.svc.Update(c.Request().Context(), id, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *OrganizationHandler) Delete(c echo.Context) error {
	id := c.Param("id")

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.NoContent(c)
}
