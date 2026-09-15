package http

import (
	"github.com/labstack/echo/v4"
)

// RegisterOrganizationRoutes registers all Organization routes on the given Echo group.
func RegisterOrganizationRoutes(g *echo.Group, h *OrganizationHandler) {
	g.POST("", h.Create)
	g.GET("/mine", h.ListMine) // Organizations where the authenticated user is a member
	g.GET("/members", h.ListMembers)
	g.POST("/members", h.AddMember)
	g.DELETE("/members/:userId", h.RemoveMember)
	g.GET("/:id", h.GetByID)
	g.GET("", h.List)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
