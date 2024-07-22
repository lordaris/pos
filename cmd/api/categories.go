package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
)

/*
* Categories are intended to be used as departments of the store, such as electronics,
* fruits and vegetables, clothing, snack, and so on. This might help to organize products in different ways,
* like with inventory management, and allow different heads of departments to place orders, check the inventory,
* only for its categories.
* */

func (app *application) createCategory(c *gin.Context) {
	var input struct {
		Name string `json:"name"`
	}

	category := &data.Category{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.Name = input.Name

	// Check if a category with the same name exists
	categoriesCollection := app.Collection(data.CollectionCategory)
	var existingCategory data.Category
	err := categoriesCollection.FindOne(context.TODO(), bson.M{"name": category.Name}).Decode(&existingCategory)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Category already exists"})
		return
	}

	result, err := categoriesCollection.InsertOne(context.TODO(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create categories"})
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}
