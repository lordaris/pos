package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
)

func (app *application) createAuthenticationToken(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usersCollection := app.Collection(data.CollectionUser)
	var user data.User
	err := usersCollection.FindOne(context.TODO(), bson.M{"username": input.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No users with that username"})
		return
	}

	match, err := user.MatchesPassword(input.Password)
	if err != nil || !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Set the duration of the token
	token, err := data.GenerateToken(time.Hour)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = usersCollection.UpdateOne(context.TODO(), bson.M{"username": user.Username}, bson.M{"$push": bson.M{"tokens": token}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"authentication_token": gin.H{"token": token.Plaintext, "expiry": token.Expiry}})
}
