package auth

import (
    "strings"

    "github.com/dgrijalva/jwt-go"
)

// User entry in the database
type UserEntity struct {
    Email        string `json:"email,omitempty"`
    PasswordHash string `json:"passwordHash,omitempty"`
    Roles        string `json:"roles"`
    LastSignIn   string `json:"lastSignIn,omitempty"`
}

func (u UserEntity) GetRoles() []string {
    cleanedRoles := make([]string, 0)
    for _, v := range strings.Split(u.Roles, ",") {
        trimmedRole := strings.Trim(v, "\t \n")
        if trimmedRole != "" {
            cleanedRoles = append(cleanedRoles, trimmedRole)
        }
    }
    return cleanedRoles
}

//JwtClaims extends the standard set of claims with roles
type JwtClaims struct {
    jwt.StandardClaims
    Roles []string `json:"roles,omitempty"`
}
