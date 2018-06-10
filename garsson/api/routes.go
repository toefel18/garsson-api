package api

func (s *Server) configureRoutes() {
    s.router.POST("/api/v1/login", s.login())

	authenticated := s.router.Group("/api")
	authenticated.Use(s.authenticate())

	v1 := authenticated.Group("/v1")
	v1.GET("/hello", s.handleHello())
	v1.GET("/db", s.databaseVersion(), s.requireRole("sjonnie"))
	v1.GET("/products", s.handleProducts())
	v1.GET("/orders", s.handleOrders())
	v1.GET("/orders/:orderId", s.handleOrder())
}