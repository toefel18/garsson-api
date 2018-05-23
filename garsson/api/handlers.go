package api

import (
    "fmt"

    "github.com/labstack/echo"
    "github.com/toefel18/garsson-api/garsson/auth"
    "github.com/toefel18/garsson-api/garsson/db/migration"
    "github.com/toefel18/garsson-api/garsson/log"

    "net/http"
)

func (s *Server) notFound() echo.HandlerFunc {
    return func(c echo.Context) error {
        return c.JSON(http.StatusNotFound, GenericResponse{Code: http.StatusNotFound, Message: "not found"})
    }
}

func (s *Server) methodNotAllowed() echo.HandlerFunc {
    return func(c echo.Context) error {
        return c.JSON(http.StatusMethodNotAllowed, GenericResponse{Code: http.StatusMethodNotAllowed, Message: "method not allowed"})
    }
}

func (s *Server) handleHello() echo.HandlerFunc {
    return func(c echo.Context) error {
        return c.JSON(http.StatusOK, "World")
    }
}

func (s *Server) login() echo.HandlerFunc {
    type LoginRequest struct {
        Email    string `json:"email" form:"email" query:"email"`
        Password string `json:"password" form:"password" query:"password"`
    }

    return func(c echo.Context) error {
        loginRequest := new(LoginRequest)
        if err := c.Bind(loginRequest); err == echo.ErrUnsupportedMediaType {
            return c.JSON(http.StatusUnsupportedMediaType, GenericResponse{Code: http.StatusUnsupportedMediaType, Message: "unsupported media type, use: application/json, application/xml or application/x-www-form-urlencoded"})
        } else if err != nil {
            return c.JSON(http.StatusBadRequest, GenericResponse{Code: http.StatusBadRequest, Message: err.Error()})
        }

        if jwt, user, err := auth.Authenticate(s.dao.NewSession(), loginRequest.Email, loginRequest.Password, []byte("dummy-for-now")); err == nil {
            log.WithField("email", user.Email).WithField("roles", user.Roles).Info("user authenticated")
            c.Response().Header().Add("Authorization", fmt.Sprintf("Bearer %v", jwt))
            return c.JSON(http.StatusOK, GenericResponse{Code: http.StatusOK, Message: "login success", Data: user})
        } else {
            log.WithField("email", loginRequest.Email).WithError(err).Warn("failed to authenticate")
            return c.JSON(http.StatusUnauthorized, GenericResponse{Code: http.StatusUnauthorized, Message: err.Error()})
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
