package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) createUser(c *gin.Context) {
	// Define a struct to hold the data from the request body
	var input struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Password string `json:"password"`
		RoleID   string `json:"role_id"`
	}

	// Create a new User struct pointer to hold user information
	user := &data.User{}

	// Bind the JSON body from the request to the `input` struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Transform the RoleID from string to ObjectID
	roleObjectID, err := primitive.ObjectIDFromHex(input.RoleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Verify if RoleID exists in the roles collection
	rolesCollection := app.Collection(data.CollectionRole)
	var role data.Role
	err = rolesCollection.FindOne(context.TODO(), bson.M{"_id": roleObjectID}).Decode(&role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found"})
		return
	}

	// Check if a user with the same username already exists
	usersCollection := app.Collection(data.CollectionUser)
	var existingUser data.User
	err = usersCollection.FindOne(context.TODO(), bson.M{"username": input.Username}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Assign data from input to the user structure
	user.Name = input.Name
	user.Username = input.Username
	user.RoleID = roleObjectID
	user.Created = time.Now()
	user.Tokens = []primitive.ObjectID{}

	// Set the user's password securely using the `SetPassword` method of the User struct (pointer)
	err = user.SetPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Insert the new user document into the user collection
	result, err := usersCollection.InsertOne(context.TODO(), user)
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
		Name     *string `json:"name"`
		Username *string `json:"username"`
		Password *string `json:"password"` // Optional field for password update
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
	usersCollection := app.Collection(data.CollectionUser)
	var existingUser data.User
	err = usersCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verify that the new username doesn't exist, and if so, return an error.
	if updatedUser.Username != nil && *updatedUser.Username != existingUser.Username {
		var existingUsername data.User
		err = usersCollection.FindOne(context.TODO(), bson.M{"username": *updatedUser.Username}).Decode(&existingUsername)
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		} else if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		existingUser.Username = *updatedUser.Username
	}

	// Update user fields only if provided
	if updatedUser.Name != nil {
		existingUser.Name = *updatedUser.Name
	}

	if updatedUser.Password != nil {
		err = existingUser.SetPassword(*updatedUser.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
			return
		}
	}

	// Update the user document in the database
	_, err = usersCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{"$set": existingUser})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (app *application) updateUserRole(c *gin.Context) {
	userID := c.Param("id")

	// Convert the ID to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	usersCollection := app.Collection(data.CollectionUser)
	var existingUser data.User
	err = usersCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&existingUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var updateRole struct {
		RoleID string `json:"role_id"`
	}

	if err := c.BindJSON(&updateRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	roleObjectID, err := primitive.ObjectIDFromHex(updateRole.RoleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	rolesCollection := app.Collection(data.CollectionRole)
	var existingRole data.Role
	err = rolesCollection.FindOne(context.TODO(), bson.M{"_id": roleObjectID}).Decode(&existingRole)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = usersCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{"$set": bson.M{"role_id": roleObjectID}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

func (app *application) getUser(c *gin.Context) {
	userID := c.Param("id")

	type userSearch struct {
		Name     string             `bson:"name"`
		Username string             `bson:"username"`
		RoleID   primitive.ObjectID `bson:"role_id"`
		Created  time.Time          `bson:"created"`
	}

	// Convert id string to ObjectId
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by ObjectId
	var user userSearch
	usersCollection := app.Collection(data.CollectionUser)
	filter := bson.D{{"_id", id}}
	err = usersCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return user
	c.JSON(http.StatusOK, user)
}

// Get all users with the same role.
func (app *application) getUsersByRole(c *gin.Context) {
	roleID := c.Param("role")

	type usersSearch struct {
		Name     string             `bson:"name"`
		Username string             `bson:"username"`
		RoleID   primitive.ObjectID `bson:"role_id"`
		Created  time.Time          `bson:"created"`
	}

	// Convert role string into ObjectID
	id, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	usersCollection := app.Collection(data.CollectionUser)
	filter := bson.D{{"role_id", id}}

	cursor, err := usersCollection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var results []usersSearch
	if err = cursor.All(context.TODO(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No users found for the given role"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func (app *application) deleteUser(c *gin.Context) {
	userID := c.Param("id")

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	usersCollection := app.Collection(data.CollectionUser)
	filter := bson.D{{"_id", id}}
	result, err := usersCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
