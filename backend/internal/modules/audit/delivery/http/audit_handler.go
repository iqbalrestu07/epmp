package http

import (
	"strconv"

	"github.com/epmp/backend/internal/modules/audit/service"
	mw "github.com/epmp/backend/internal/pkg/middleware"
	"github.com/epmp/backend/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

// AuditHandler handles HTTP requests for the audit trail.
type AuditHandler struct {
	svc *service.AuditService
}

// NewAuditHandler creates a new AuditHandler.
func NewAuditHandler(svc *service.AuditService) *AuditHandler {
	return &AuditHandler{svc: svc}
}

func (h *AuditHandler) List(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return response.BadRequest(c, "X-Organization-ID header is required")
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage == 0 {
		perPage = 20
	}

	result, err := h.svc.List(c.Request().Context(), page, perPage,
		c.QueryParam("search"), c.QueryParam("module"), c.QueryParam("action"), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

// RegisterAuditRoutes registers all Audit routes on the given Echo group.
func RegisterAuditRoutes(g *echo.Group, h *AuditHandler) {
	g.GET("", h.List)
}
