package auth

import (
    "encoding/hex"
    "errors"
    "fmt"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/gocraft/dbr"
    "github.com/satori/go.uuid"
    "golang.org/x/crypto/sha3"
)

var (
    // ErrUserNotFound indicates that the user does not exist in the database
    ErrUserNotFound    = errors.New("user not found")
    // ErrInvalidPassword indicates that the passwords was incorrect
    ErrInvalidPassword = errors.New("invalid password")
)

const (
    TokenValidity           = time.Hour * 8
    TokenGenerationErrorFmt = "could not generate token: %v"
)

func Authenticate(sess dbr.SessionRunner, email, password string, signingSecret []byte) (jwt string, user UserEntity, err error) {
    if user, err = QueryUserEntity(sess, email); err == dbr.ErrNotFound {
        return "", UserEntity{}, ErrUserNotFound
    } else if hashPassword(password) != user.PasswordHash {
        return "", UserEntity{}, ErrInvalidPassword
    } else if jwt, err = createToken(user, signingSecret); err != nil {
        return "", UserEntity{}, fmt.Errorf(TokenGenerationErrorFmt, err.Error())
    } else {
        return
    }
}

func hashPassword(password string) string {
    blaat := sha3.Sum512([]byte(password))
    return hex.EncodeToString(blaat[:])
}

func createToken(user UserEntity, signingSecret []byte) (string, error) {
    claims := JwtClaims{
       StandardClaims: jwt.StandardClaims{
           Id:        uuid.NewV4().String(),
           Issuer:    "garsson-api",
           IssuedAt:  time.Now().Unix(),
           NotBefore: time.Now().Unix(),
           ExpiresAt: time.Now().Add(TokenValidity).Unix(),
           Audience:  "dhl-internal",
           Subject:   user.Email,
       },
       Roles: user.GetRoles(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(signingSecret)
}
