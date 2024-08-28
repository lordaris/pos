package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
)

func (app *application) createRoles(c *gin.Context) {
	var roles struct {
		Name string `json:"name"`
	}

	// Bind the JSON body from the request to the `roles` slice
	if err := c.ShouldBindJSON(&roles); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rolesCollection := app.Collection(data.CollectionRole)
	var existingRole data.Role
	err := rolesCollection.FindOne(context.TODO(), bson.M{"name": roles.Name}).Decode(&existingRole)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Role already exists", "Role": existingRole.Name})
		return
	}

	roleData := &data.Role{
		Name: roles.Name,
	}

	result, err := rolesCollection.InsertOne(context.TODO(), roleData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ID": result.InsertedID, "Name": roleData.Name})
}

func (app *application) getAllRoles(c *gin.Context) {
	rolesCollection := app.Collection(data.CollectionRole)
	cursor, err := rolesCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching roles"})
		return
	}

	var roles []data.Role

	for cursor.Next(context.TODO()) {
		var role data.Role
		if err := cursor.Decode(&role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding roles"})
			return
		}
		roles = append(roles, role)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}
