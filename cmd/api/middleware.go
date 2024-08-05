package main

import (
	"crypto/sha256"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add the "Vary: Authorization" header to the response.
		c.Writer.Header().Add("Vary", "Authorization")

		// Retrieve the value of the Authorization header from the request.
		authorizationHeader := c.GetHeader("Authorization")

		// Split the Authorization header into its constituent parts.
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
			c.Abort()
			return
		}

		// Extract the actual authentication token from the header parts.
		token := headerParts[1]

		// Validate the token format.
		if err := data.ValidateTokenPlaintext(token); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Hash the token plaintext.
		tokenHash := sha256.Sum256([]byte(token))

		// Get token document from the database.
		usersCollection := app.Collection(data.CollectionUser)

		var user data.User
		err := usersCollection.FindOne(c, bson.M{"tokens.hash": tokenHash[:]}).Decode(&user)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
			c.Abort()
			return
		}

		// Remove expired tokens
		// Get the actual time
		currentTime := time.Now()
		// Filter the valid tokens and create a new list of validTokens.
		validTokens := user.Tokens[:0]
		for _, token := range user.Tokens {
			if token.Expiry.After(currentTime) {
				validTokens = append(validTokens, token)
			}
		}

		// If the length of valid tokens is different from the length of the original list of tokens (user.Tokens) is different, there were expired tokens, so the list of tokens is updated to show the new list of tokens.
		if len(validTokens) != len(user.Tokens) {
			_, err := usersCollection.UpdateOne(
				c,
				bson.M{"_id": user.ID},
				bson.M{"$set": bson.M{"tokens": validTokens}},
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user tokens"})
				c.Abort()
				return
			}
		}

		// Add the user to the context.
		app.contextSetUser(c, &user)

		// Call the next handler in the chain.
		c.Next()
	}
}
