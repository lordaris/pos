package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (app *application) allowRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := app.contextGetUser(c)
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			c.Abort()
			return
		}

		rolesCollection := app.Collection(data.CollectionRole)
		var roleInfo data.Role
		err := rolesCollection.FindOne(c, bson.M{"_id": user.RoleID}).Decode(&roleInfo)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid role ID "})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve role"})
			c.Abort()
			return
		}

		allowed := false
		for _, role := range allowedRoles {
			if roleInfo.Name == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You don't have permission to use this "})
			c.Abort()
			return
		}
		c.Next()
	}
}

// CustomResponseWriter wraps the gin.ResponseWriter to capture the response body
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body in addition to writing it to the original ResponseWriter
func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// loggerMiddleware is a Gin middleware function that logs details about each request and response
func (app *application) loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		var requestBody string

		// For POST, PUT, and PATCH requests, read and store the request body
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			// Read the raw body data
			rawBody, err := io.ReadAll(c.Request.Body)
			if err == nil && len(rawBody) > 0 {
				// Convert raw body to string
				requestBody = string(rawBody)
				// Restore the request body for further processing
				c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))
				// Hide passwords
				requestBody = app.hidePassword(requestBody)
			}
		}

		// Create a custom ResponseWriter to capture the response body
		responseBody := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseBody

		// Process the request
		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		var userID primitive.ObjectID
		var username string

		// If a user is authenticated, get their details
		if user := app.contextGetUser(c); user != nil {
			userID = user.ID
			username = user.Username
		}

		logEntry := data.LogEntry{
			Timestamp:    startTime,
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			Status:       c.Writer.Status(),
			Latency:      latency.String(),
			UserID:       userID,
			Username:     username,
			RequestBody:  requestBody,
			ResponseBody: responseBody.body.String(),
		}

		// Add method-specific information to the log entry
		switch c.Request.Method {
		case "GET":
			// For GET requests, log the query parameters
			logEntry.QueryParams = c.Request.URL.RawQuery
		case "DELETE":
			// For DELETE requests, log the resource ID if available
			resourceID := c.Param("id")
			if resourceID != "" {
				logEntry.ResourceID = resourceID
			}
		}

		// Convert the log entry to JSON
		logJSON, err := json.Marshal(logEntry)
		if err == nil {
			// Print the JSON log entry
			fmt.Println(string(logJSON))
		}

		// Save the log entry to the database
		app.saveLogToDatabase(logEntry)
	}
}

func (app *application) saveLogToDatabase(logEntry data.LogEntry) {
	collection := app.Collection("logs")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, logEntry)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	log.Println("Log saved")
}

func (app *application) hidePassword(jsonStr string) string {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return jsonStr
	}

	if _, ok := data["password"]; ok {
		data["password"] = "Hidden password"
	}

	encryptedJSON, err := json.Marshal(data)
	if err != nil {
		return jsonStr
	}

	return string(encryptedJSON)
}
