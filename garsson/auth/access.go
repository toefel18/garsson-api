package auth

import (
    "github.com/gocraft/dbr"
    "github.com/toefel18/garsson-api/garsson/db"
)

// QueryUserEntity returns the user as stored in the db, or dbr.ErrNotFound if no user was found
// with that email.
func QueryUserEntity(session dbr.SessionRunner, email string) (user UserEntity, err error) {
    err = session.
        Select("*").
        From(db.UserAccountTable).
        Where("email = ?", email).
        LoadOne(&user)
    return
}

func UpdateLastSignInToNow(session dbr.SessionRunner, email string) error {
    // TODO implement: update LastSignIn to now!
    return nil
}