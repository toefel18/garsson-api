package api

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/toefel18/garsson-api/garsson/db"
)

func Publish(dao *db.Dao) {
	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Secure())
	router.Logger.Fatal(router.Start(":8080"))
}
