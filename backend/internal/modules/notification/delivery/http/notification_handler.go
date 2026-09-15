package http

import (
	"strconv"

	"github.com/epmp/backend/internal/modules/notification/service"
	mw "github.com/epmp/backend/internal/pkg/middleware"
	"github.com/epmp/backend/internal/pkg/response"

	"github.com/labstack/echo/v4"
)

// NotificationHandler handles HTTP requests for in-app notifications.
type NotificationHandler struct {
	svc *service.NotificationService
}

// NewNotificationHandler creates a new NotificationHandler.
func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) List(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return response.BadRequest(c, "X-Organization-ID header is required")
	}
	userID := mw.GetUserID(c)

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}
	perPage, _ := strconv.Atoi(c.QueryParam("per_page"))
	if perPage == 0 {
		perPage = 20
	}

	result, err := h.svc.List(c.Request().Context(), userID, orgID, page, perPage)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.OK(c, result)
}

func (h *NotificationHandler) MarkRead(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return response.BadRequest(c, "X-Organization-ID header is required")
	}

	if err := h.svc.MarkRead(c.Request().Context(), c.Param("id"), orgID); err != nil {
		return response.NotFound(c, "Notification not found")
	}

	return response.NoContent(c)
}

func (h *NotificationHandler) MarkAllRead(c echo.Context) error {
	orgID := mw.GetOrgID(c)
	if orgID == "" {
		return response.BadRequest(c, "X-Organization-ID header is required")
	}

	if err := h.svc.MarkAllRead(c.Request().Context(), mw.GetUserID(c), orgID); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.NoContent(c)
}

// RegisterNotificationRoutes registers all Notification routes on the given Echo group.
func RegisterNotificationRoutes(g *echo.Group, h *NotificationHandler) {
	g.GET("", h.List)
	g.POST("/:id/read", h.MarkRead)
	g.POST("/read-all", h.MarkAllRead)
}
