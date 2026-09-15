package http

import (
	"strconv"

	"github.com/epmp/backend/internal/modules/reservation/dto"
	"github.com/epmp/backend/internal/modules/reservation/service"
	"github.com/epmp/backend/internal/pkg/errs"
	mw "github.com/epmp/backend/internal/pkg/middleware"
	"github.com/epmp/backend/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

// ReservationHandler handles HTTP requests for Reservation resources.
type ReservationHandler struct {
	svc *service.ReservationService
}

// NewReservationHandler creates a new ReservationHandler.
func NewReservationHandler(svc *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{svc: svc}
}

func (h *ReservationHandler) Create(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return response.BadRequest(c, "X-Organization-ID header is required")
	}

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	result, err := h.svc.Create(c.Request().Context(), orgID, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, result)
}

func (h *ReservationHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	result, err := h.svc.GetByID(c.Request().Context(), id, orgID)
	if err != nil {
		return response.NotFound(c, "Reservation not found")
	}

	return response.OK(c, result)
}

func (h *ReservationHandler) List(c echo.Context) error {
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

func (h *ReservationHandler) Update(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	var req dto.UpdateReservationRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	result, err := h.svc.Update(c.Request().Context(), id, orgID, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *ReservationHandler) Delete(c echo.Context) error {
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
