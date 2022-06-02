package http

import (
	"github.com/labstack/echo/v4"

	"seeder/internal/domain"
)

//MapNodeRoutes Map node routes
func MapNodeRoutes(nodeGroup *echo.Group, h domain.NodeHandlers) {
	nodeGroup.GET("", h.GetNodesList())
	nodeGroup.POST("", h.AuthNode())
}
