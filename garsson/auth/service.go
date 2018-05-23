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
    ErrUserNotFound = errors.New("user not found")
    // ErrInvalidPassword indicates that the passwords was incorrect
    ErrInvalidPassword = errors.New("invalid password")
)

const (
    // TokenValidity is the time in which the token is valid
    TokenValidity = time.Hour * 8
    // TokenGenerationErrorFmt contains the error format when token generation does not work
    TokenGenerationErrorFmt = "could not generate token: %v"
)

// Authenticate checks if a user has the right credentials and provides a JWT token, along with the users entity
func Authenticate(sess dbr.SessionRunner, email, password string, signingSecret []byte) (string, UserEntity, error) {
    if user, err := QueryUserEntity(sess, email); err == dbr.ErrNotFound {
        return "", UserEntity{}, ErrUserNotFound
    } else if hashPassword(password) != user.PasswordHash {
        return "", UserEntity{}, ErrInvalidPassword
    } else if jwt, err := createToken(user, signingSecret); err != nil {
        return "", UserEntity{}, fmt.Errorf(TokenGenerationErrorFmt, err.Error())
    } else {
        user.PasswordHash = "" // no need to expose!
        return jwt, user, nil
    }
}

func hashPassword(password string) string {
    rawHash := sha3.Sum512([]byte(password))
    return hex.EncodeToString(rawHash[:])
}

func createToken(user UserEntity, signingSecret []byte) (string, error) {
    claims := JwtClaims{
        StandardClaims: jwt.StandardClaims{
            Id:        uuid.NewV4().String(),
            Issuer:    "garsson-api",
            IssuedAt:  time.Now().Unix(),
            NotBefore: time.Now().Add(-2 * time.Minute).Unix(), // -2 minutes to allow for clock drift
            ExpiresAt: time.Now().Add(TokenValidity).Unix(),
            Audience:  "garsson-api-users",
            Subject:   user.Email,
        },
        Roles: user.GetRoles(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(signingSecret)
}
