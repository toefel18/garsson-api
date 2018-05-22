package api

import (
	"github.com/labstack/echo/middleware"
)

func (s *Server) configureRoutes() {
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recover())
	s.router.Use(middleware.Secure())

	s.router.POST("/login", s.login())
	s.router.GET("/hello", s.handleHello())
	s.router.GET("/db", s.databaseVersion())
}