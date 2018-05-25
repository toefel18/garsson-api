package api

func (s *Server) configureRoutes() {
    // unauthenticated route
	s.router.POST("/v1/login", s.login())

	authenticated := s.router.Group("")
	authenticated.Use(s.authenticate())

	v1 := authenticated.Group("/v1")
	v1.GET("/hello", s.handleHello())
	v1.GET("/db", s.databaseVersion(), s.requireRole("sjonnie"))
}