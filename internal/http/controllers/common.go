package controllers

import (
	"github.com/labstack/echo/v4"
	"vehicles/internal/domain"
	m "vehicles/internal/http/middlewares"
)

func parseUser(c echo.Context) *domain.User {
	if c.Get("user") == nil {
		user := m.GenerateMetadata(c, nil)
		return user
	}
	return c.Get("user").(*domain.User)
}
