package api

import (
    "fmt"

    "github.com/labstack/echo"
    "github.com/toefel18/garsson-api/garsson/auth"
    "github.com/toefel18/garsson-api/garsson/db/migration"

    "net/http"
)

func (s *Server) handleHello() echo.HandlerFunc {
    return func(c echo.Context) error {
        return c.JSON(http.StatusOK, "World")
    }
}

func (s *Server) login() echo.HandlerFunc {
    type LoginRequest struct {
        Email    string
        Password string
    }

    return func(c echo.Context) error {
        loginRequest := new(LoginRequest)
        if err := c.Bind(loginRequest); err != nil {
            return c.JSON(http.StatusBadRequest, GenericResponse{Code: http.StatusBadRequest, Message: err.Error()})
        }
        if jwt, user, err := auth.Authenticate(s.dao.NewSession(), loginRequest.Email, loginRequest.Password, []byte("dummy-for-now")); err == nil {
            c.Response().Header().Add("Authorization", fmt.Sprintf("Bearer %v", jwt))
            c.JSON(http.StatusOK, )
        }

    }
}

func (s *Server) databaseVersion() echo.HandlerFunc {
    return func(c echo.Context) error {
        if versions, err := migration.FetchDbVersion(s.dao.NewSession()); err != nil {
            return c.JSON(http.StatusInternalServerError, GenericResponse{Code: http.StatusInternalServerError, Message: err.Error()})
        } else {
            return c.JSON(http.StatusOK, versions)
        }
    }
}
