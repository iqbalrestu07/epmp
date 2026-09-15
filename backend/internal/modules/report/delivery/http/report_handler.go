package http

import (
	"time"

	"github.com/epmp/backend/internal/modules/report/service"
	mw "github.com/epmp/backend/internal/pkg/middleware"
	"github.com/epmp/backend/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

// ReportHandler handles HTTP requests for read-only reports.
type ReportHandler struct {
	svc *service.ReportService
}

// NewReportHandler creates a new ReportHandler.
func NewReportHandler(svc *service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) GetOccupancy(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return response.BadRequest(c, "X-Organization-ID header is required")
	}

	result, err := h.svc.GetOccupancyReport(c.Request().Context(), orgID, c.QueryParam("property_id"), c.QueryParam("building_id"))
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *ReportHandler) GetRevenue(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return response.BadRequest(c, "X-Organization-ID header is required")
	}

	now := time.Now()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -5, 0)
	to := from.AddDate(0, 6, 0)

	if v := c.QueryParam("from"); v != "" {
		parsed, err := time.Parse("2006-01-02", v)
		if err != nil {
			return response.BadRequest(c, "invalid 'from' date, expected YYYY-MM-DD")
		}
		from = parsed
	}
	if v := c.QueryParam("to"); v != "" {
		parsed, err := time.Parse("2006-01-02", v)
		if err != nil {
			return response.BadRequest(c, "invalid 'to' date, expected YYYY-MM-DD")
		}
		to = parsed
	}
	if !to.After(from) {
		return response.BadRequest(c, "'to' must be after 'from'")
	}

	result, err := h.svc.GetRevenueReport(c.Request().Context(), orgID, from, to)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *ReportHandler) GetArAging(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return response.BadRequest(c, "X-Organization-ID header is required")
	}

	result, err := h.svc.GetArAgingReport(c.Request().Context(), orgID)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

// RegisterReportRoutes registers all Report routes on the given Echo group.
func RegisterReportRoutes(g *echo.Group, h *ReportHandler) {
	g.GET("/occupancy", h.GetOccupancy)
	g.GET("/revenue", h.GetRevenue)
	g.GET("/ar-aging", h.GetArAging)
}
