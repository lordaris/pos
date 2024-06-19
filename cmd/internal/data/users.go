package data

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `bson:"name"`
	Username       string             `bson:"username"`
	HashedPassword []byte             `bson:"password"` // Store hashed password only
	Created        time.Time          `bson:"created"`
	RoleID         primitive.ObjectID `bson:"role_id"`
}

func (u *User) SetPassword(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	u.HashedPassword = hash
	return nil
}

func (u *User) MatchesPassword(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
