package auth

import (
    "strings"

    "github.com/dgrijalva/jwt-go"
    "github.com/kubernetes/kubernetes/pkg/util/slice"
)

const (
    //RoleAdmin is the role name of administrators
    RoleAdmin = "admin"
)

// User entry in the database
type UserEntity struct {
    // Email acts as the main user identifier
    Email string `json:"email,omitempty"`
    // PasswordHash is the sha3 hash, hex encoded
    PasswordHash string `json:"passwordHash,omitempty"`
    // Roles is a csv string with all the roles owned by the user
    Roles string `json:"roles"`
    // LastSignIn contains the timestamp of last sign-in.
    LastSignIn string `json:"lastSignIn,omitempty"`
}

// GetRoles returns all the roles of the user as an array
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

// UserFromJwt holds the fields found in a JWT
type UserFromJwt struct {
    // Email is the user id, equal to the sub field of the claims
    Email   string
    // Roles that are listed in the jwt
    Roles  []string
    // Claims contains a reference to the other claims found in the JWT
    Claims *JwtClaims
}

// HasRole checks if the user owns a given role. Returns true for every role if the user also has the role admin!
// also returns the role that granted access (either the role, or admin)
func (u UserFromJwt) HasRole(role string) (bool, string) {
    if slice.ContainsString(u.Roles, RoleAdmin, nil) {
        return true, RoleAdmin
    } else if slice.ContainsString(u.Roles, role, nil) {
        return true, role
    }
    return false, ""
}

//JwtClaims extends the standard set of claims with roles
type JwtClaims struct {
    jwt.StandardClaims
    Roles []string `json:"roles,omitempty"`
}
