package api

import (
    "github.com/labstack/echo"
    "github.com/toefel18/garsson-api/garsson/db/migration"

    "net/http"
)

func respondWithError(c echo.Context, code int, err error) error {
    c.JSON(code, GenericResponse{Code: code, Message: err.Error()})
    return err
}

func (s *Server) handleHello() echo.HandlerFunc {
    return func(c echo.Context) error {
        return c.JSON(http.StatusOK, "World")
    }
}

func (s *Server) databaseVersion() echo.HandlerFunc {
    return func(c echo.Context) error {
        if versions, err := migration.FetchDbVersion(s.dao.NewSession()); err != nil {
            return respondWithError(c, http.StatusInternalServerError, err)
        } else {
            return c.JSON(http.StatusOK, versions)
        }
    }
}
