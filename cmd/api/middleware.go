package main

import (
	"crypto/sha256"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
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
		tokensCollection := app.config.db.mongoClient.Database("pos").Collection("tokens")
		usersCollection := app.config.db.mongoClient.Database("pos").Collection("user")

		var tokenDoc data.Token
		err := tokensCollection.FindOne(c, bson.M{"hash": tokenHash[:]}).Decode(&tokenDoc)
		if err != nil {
			// If no token is found, respond with an error.
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
			c.Abort()
			return
		}

		// Check if the token is expired.
		if tokenDoc.Expiry.Before(time.Now()) {
			// Remove the expired token from the database.
			_, err := tokensCollection.DeleteOne(c, bson.M{"_id": tokenDoc.ID})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove expired token"})
				c.Abort()
				return
			}

			// Remove the token reference from the user document.
			usersCollection := app.config.db.mongoClient.Database("pos").Collection("user")
			_, err = usersCollection.UpdateOne(
				c,
				bson.M{"username": tokenDoc.UserName},
				bson.M{"$pull": bson.M{"tokens": tokenDoc.ID}},
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user tokens"})
				c.Abort()
				return
			}

			c.JSON(http.StatusUnauthorized, gin.H{"error": "token has expired and has been removed, please log in again"})
			c.Abort()
			return
		}
		var user data.User
		err = usersCollection.FindOne(c, bson.M{"username": tokenDoc.UserName}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		// Add the user to the context.
		app.contextSetUser(c, &user)

		// Call the next handler in the chain.
		c.Next()
	}
}
