package api

import (
	"github.com/labstack/echo"
	"net/http"
)

func (s *Server) handleHello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, "World")
	}
}
