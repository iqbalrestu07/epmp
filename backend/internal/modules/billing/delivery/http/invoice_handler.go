package http

import (
	"strconv"

	"github.com/epmp/backend/internal/modules/billing/dto"
	"github.com/epmp/backend/internal/modules/billing/service"
	"github.com/epmp/backend/internal/pkg/errs"
	mw "github.com/epmp/backend/internal/pkg/middleware"
	"github.com/epmp/backend/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

// InvoiceHandler handles HTTP requests for Invoice resources.
type InvoiceHandler struct {
	svc *service.InvoiceService
}

// NewInvoiceHandler creates a new InvoiceHandler.
func NewInvoiceHandler(svc *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{svc: svc}
}

func (h *InvoiceHandler) Create(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return response.BadRequest(c, "X-Organization-ID header is required")
	}

	var req dto.CreateInvoiceRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	result, err := h.svc.Create(c.Request().Context(), orgID, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, result)
}

func (h *InvoiceHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	result, err := h.svc.GetByID(c.Request().Context(), id, orgID)
	if err != nil {
		return response.NotFound(c, "Invoice not found")
	}

	return response.OK(c, result)
}

func (h *InvoiceHandler) List(c echo.Context) error {
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

func (h *InvoiceHandler) Update(c echo.Context) error {
	id := c.Param("id")
	orgID := mw.GetOrgID(c)

	var req dto.UpdateInvoiceRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	result, err := h.svc.Update(c.Request().Context(), id, orgID, &req)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *InvoiceHandler) Delete(c echo.Context) error {
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
