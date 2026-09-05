package http

import (
	"github.com/epmp/backend/internal/modules/dashboard/service"
	mw "github.com/epmp/backend/internal/pkg/middleware"
	"github.com/epmp/backend/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) GetSummary(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return response.BadRequest(c, "X-Organization-ID header is required")
	}

	result, err := h.svc.GetSummary(c.Request().Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func RegisterDashboardRoutes(g *echo.Group, h *DashboardHandler) {
	g.GET("/summary", h.GetSummary)
}
