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

	_, err = rolesCollection.InsertOne(context.TODO(), roleData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": roleData})
}
