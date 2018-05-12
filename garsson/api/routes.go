package api

import (
	"github.com/labstack/echo/middleware"
)

func (s *Server) configureRoutes() {
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Secure())
	s.router.GET("/hello", s.handleHello())
}