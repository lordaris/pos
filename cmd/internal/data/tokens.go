package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Plaintext string             `bson:"token"`
	Hash      []byte             `bson:"hash"`
	UserName  string             `bson:"username"`
	Expiry    time.Time          `bson:"expiry"`
}

func GenerateToken(userName string, ttl time.Duration) (*Token, error) {
	token := &Token{
		ID:       primitive.NewObjectID(),
		UserName: userName,
		Expiry:   time.Now().Add(ttl),
	}

	// initialize a zero-valued byte slice with a lenght of 16 bytes
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	/*
			* This is not 16 characters long but has an underlying entropy of 16 bytes -of randomness-.
		  * the length depends on how those 16 random bytes are encoded to create a string, which in this case is a base-32 string,
		  * which results in a string with 26 characters. Using base-16 would result in a string with 32 characters.
			* */
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(tokenPlaintext string) error {
	if tokenPlaintext == "" {
		return errors.New("token must be provided")
	}
	if len(tokenPlaintext) != 26 {
		return errors.New("token must be 26 characters long")
	}
	return nil
}
