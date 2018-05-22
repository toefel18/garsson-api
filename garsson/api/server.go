package api

import (
	"github.com/labstack/echo"
	"github.com/toefel18/garsson-api/garsson/db"
	"github.com/toefel18/garsson-api/garsson/log"
)

// Implementation inspired by https://medium.com/@matryer/how-i-write-go-http-services-after-seven-years-37c208122831

type Server struct {
	router *echo.Echo
	dao    *db.Dao
}

func NewServer(dao *db.Dao) *Server  {
	return &Server {
		router: echo.New(),
		dao: dao,
	}
}

func (s *Server) Start() {
	s.configureRoutes()
	log.Fatal(s.router.Start(":8080"))
}