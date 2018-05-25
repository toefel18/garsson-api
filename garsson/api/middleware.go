package api

import (
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "github.com/toefel18/garsson-api/garsson/auth"
    "github.com/toefel18/garsson-api/garsson/log"
)

const (
    // bearerPrefix that every Authorization header should start with
    bearerPrefix = "Bearer "
    // bearerPrefixLen contains the length of 'Bearer '
    bearerPrefixLen = 7
    // AuthenticatedUserKey is the key to lookup the authenticated user via echo.Context.Get()
    AuthenticatedUserKey = "authenticated_user"
    // GrantingRoleKey is the key to lookup the role that granted the user access
    GrantingRoleKey = "granting_role"
    // MissingRoleKey is the key to lookup the role that the user is missing to gain access
    MissingRoleKey = "missing_role"
)

func (s *Server) configureMiddleware() {
    s.router.Use(s.loggingMiddleware())
    //s.router.Use(middleware.Logger())
    s.router.Use(middleware.Recover())
    s.router.Use(middleware.Secure())

    echo.NotFoundHandler = s.notFound()
    echo.MethodNotAllowedHandler = s.methodNotAllowed()
}

func (s *Server) loggingMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) (err error) {
            req := c.Request()
            res := c.Response()
            start := time.Now()
            if err = next(c); err != nil {
                c.Error(err)
            }
            stop := time.Now()
            username := "not_authenticated"
            if user, err := s.getCurrentUser(c); err == nil {
                username = user.Email
            }

            logStmt := log.WithFields(log.Fields{
                "method":       req.Method,
                "uri":          req.RequestURI,
                "remoteIP":     c.RealIP(),
                "host":         req.Host,
                "status":       res.Status,
                "user":         username,
                "bytesIn":      req.Header.Get(echo.HeaderContentLength),
                "bytesOut":     strconv.FormatInt(res.Size, 10),
                "latency":      strconv.FormatInt(int64(stop.Sub(start)), 10),
                "latencyHuman": stop.Sub(start).String(),
            })

            if err != nil {
                logStmt = logStmt.WithError(err)
            }
            missingRole := c.Get(MissingRoleKey)
            if missingRole != nil {
                logStmt = logStmt.WithField("userIsMissingRole", missingRole)
            }
            grantingRole := c.Get(GrantingRoleKey)
            if grantingRole != nil {
                logStmt = logStmt.WithField("accessGrantedByRole", grantingRole)
            }
            if req.Referer() != "" {
                logStmt = logStmt.WithField("referer", req.Referer())
            }

            if res.Status >= 200 && res.Status < 300 {
                logStmt.Info("[API]")
            } else if res.Status >= 500 {
                logStmt.Error("[API]")
            } else {
                logStmt.Warn("[API]")
            }

            return
        }
    }
}

// authenticate parses the JWT and sets a variable in the context, stops processing when not authenticated
func (s *Server) authenticate() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) (error) {
            req := c.Request()
            authHeader := req.Header.Get(echo.HeaderAuthorization)
            if authHeader == "" {
                return c.JSON(http.StatusUnauthorized, GenericResponse{Code: http.StatusUnauthorized, Message: "Authorization header not set, provide 'Authorization: Bearer <jwt>', acquire jwt via /v1/login"})
            } else if !strings.HasPrefix(authHeader, bearerPrefix) {
                return c.JSON(http.StatusUnauthorized, GenericResponse{Code: http.StatusUnauthorized, Message: "Authorization header does not start with 'Bearer '"})
            } else if user, err := auth.ValidateJWT(authHeader[bearerPrefixLen:], s.jwtSigningSecret); err != nil {
                return c.JSON(http.StatusUnauthorized, GenericResponse{Code: http.StatusUnauthorized, Message: err.Error()})
            } else {
                c.Set(AuthenticatedUserKey, user)
                return next(c)
            }
        }
    }
}

// requireRole assumes authenticate has already run, checks if the user has the required role
func (s *Server) requireRole(role string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) (error) {
            user, err := s.getCurrentUser(c)
            if err != nil {
                log.WithError(err).Warn("requireRole expects that authenticate() middleware has run, appears not!")
                return c.JSON(http.StatusUnauthorized, GenericResponse{Code: http.StatusUnauthorized, Message: "not authenticated"})
            }
            hasRole, grantingRole := user.HasRole(role)
            if !hasRole {
                log.WithFields(log.Fields{
                    "user": user.Email,
                    "userRoles": user.Roles,
                    "requiredRole": role,
                }).Warn("blocked unauthorized access")
                c.Set(MissingRoleKey, role)
                return c.JSON(http.StatusForbidden, GenericResponse{Code: http.StatusForbidden, Message: "not authorized"})
            } else {
                c.Set(GrantingRoleKey, grantingRole)
            }
            return next(c)
        }
    }
}
