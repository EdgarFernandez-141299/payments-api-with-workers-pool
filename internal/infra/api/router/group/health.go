package group

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	healthPath = "/health"
)

type HealthRoutes struct {
}

func NewHealthRoutes(
	group *echo.Group,
) *HealthRoutes {
	healthRoutes := group.Group(healthPath)
	healthRoutes.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok")
	})

	return &HealthRoutes{}
}
