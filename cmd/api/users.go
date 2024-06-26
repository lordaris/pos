package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) createUser(c *gin.Context) {
	// Define a struct to hold the data from the request body
	var input struct {
		Name     string `bson:"name"`
		Username string `bson:"username"`
		Password string `bson:"password"`
	}

	// Create a new User struct pointer to hold user information
	user := &data.User{
		Name:     input.Name,
		Username: input.Username,
	}

	collection := app.config.db.mongoClient.Database("pos").Collection("user")

	// Bind the JSON body from the request to the `input` struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a user with the same username already exists
	var existingUser data.User
	err := collection.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Set the user's password securely using the `SetPassword` method of the User struct (pointer)
	err = user.SetPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Insert the new user document into the user collection
	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

func (app *application) updateUser(c *gin.Context) {
	// Get the ID of the user to update from the request URL parameter
	userID := c.Param("id")

	// Convert the ID to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Define a struct to hold the optional updated user data
	type updateUser struct {
		Name     *string `bson:"name"`
		Username *string `bson:"username"`
		Password *string `bson:"password"` // Optional field for password update
	}

	var updatedUser updateUser

	// Bind the JSON body from the request to the `updatedUser` struct
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if at least one field is provided for update
	if updatedUser.Name == nil && updatedUser.Username == nil && updatedUser.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields provided for update"})
		return
	}

	// Get the existing user document from the database
	collection := app.config.db.mongoClient.Database("pos").Collection("user")
	var existingUser data.User
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update user fields only if provided
	if updatedUser.Name != nil {
		existingUser.Name = *updatedUser.Name
	}
	if updatedUser.Username != nil {
		existingUser.Username = *updatedUser.Username
	}

	// Update password if provided (same logic as before)
	if updatedUser.Password != nil {
		err = existingUser.SetPassword(*updatedUser.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}
	}

	// Update the user document in the database
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{"$set": existingUser})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// TODO: Delete this handler, as this is only for educational purposes.

// Implemention of the /users route that returns all of the users from the user collection
func (app *application) getUsers(c *gin.Context) {
	// Find users
	cursor, err := app.config.db.mongoClient.Database("pos").Collection("user").Find(context.TODO(), bson.D{{}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map results
	var users []bson.M
	if err = cursor.All(context.TODO(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return users
	c.JSON(http.StatusOK, users)
}

// TODO: Delete this handler as it was created only for educational purposes.

// The implementation of our /user/{id} endpoint that returns a single user based on the provided ID
func (app *application) getUserByID(c *gin.Context) {
	// Get movie ID from URL
	idStr := c.Param("id")

	// Convert id string to ObjectId
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by ObjectId
	var user bson.M
	err = app.config.db.mongoClient.Database("pos").Collection("user").FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return user
	c.JSON(http.StatusOK, user)
}

// TODO: Check for aggregations in mongodb atlas.
