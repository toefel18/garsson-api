package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (s *Server) configureRoutes() {
	s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recover())
	s.router.Use(middleware.Secure())

	echo.NotFoundHandler = s.notFound()
	echo.MethodNotAllowedHandler = s.methodNotAllowed()

	v1 := s.router.Group("/v1")

	v1.POST("/login", s.login())
	v1.GET("/hello", s.handleHello())
	v1.GET("/db", s.databaseVersion())
}