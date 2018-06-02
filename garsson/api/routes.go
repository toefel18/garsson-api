package api

import (
    "net/http"

    "github.com/labstack/echo"
)

func (s *Server) configureRoutes() {
    // unauthenticated route
    s.router.Static("/app", "/home/hestersco/Documents/projects/go/src/github.com/toefel18/garsson-api/ui")
    s.router.GET("/", func(c echo.Context) error {
        return c.Redirect(http.StatusMovedPermanently, "/app/")
    })
    s.router.GET("/app", func(c echo.Context) error {
        return c.Redirect(http.StatusMovedPermanently, "/app/")
    })

    s.router.POST("/api/v1/login", s.login())

	authenticated := s.router.Group("/api")
	authenticated.Use(s.authenticate())

	v1 := authenticated.Group("/v1")
	v1.GET("/hello", s.handleHello())
	v1.GET("/db", s.databaseVersion(), s.requireRole("sjonnie"))
}