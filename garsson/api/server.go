package api

import (
    "errors"

    "github.com/labstack/echo"
    "github.com/toefel18/garsson-api/garsson/auth"
    "github.com/toefel18/garsson-api/garsson/db"
    "github.com/toefel18/garsson-api/garsson/log"
)

// Implementation inspired by https://medium.com/@matryer/how-i-write-go-http-services-after-seven-years-37c208122831

var (
    // ErrNotAuthenticated indicates that the user is not authenticated
    ErrNotAuthenticated = errors.New("not authenticated")
)

type Server struct {
    router           *echo.Echo
    dao              *db.Dao
    jwtSigningSecret []byte
}

func NewServer(dao *db.Dao) *Server {
    return &Server{
        router:           echo.New(),
        dao:              dao,
        jwtSigningSecret: []byte("dummy-for-now"),
    }
}

func (s *Server) Start() {
    s.configureMiddleware()
    s.configureRoutes()
    log.Fatal(s.router.Start(":8080"))
}

// getCurrentUser returns the current user, if present. Requires the that the authorize middleware has run
func (s *Server) getCurrentUser(c echo.Context) (auth.UserFromJwt, error) {
    user := c.Get(AuthenticatedUserKey)
    if userObj, ok := user.(auth.UserFromJwt); ok {
        return userObj, nil

    }
    return auth.UserFromJwt{}, ErrNotAuthenticated
}
