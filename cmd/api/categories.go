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

	// Check if a category with the same name exists
	categoriesCollection := app.Collection(data.CollectionCategory)
	var existingCategory data.Category
	err := categoriesCollection.FindOne(context.TODO(), bson.M{"name": input.Name}).Decode(&existingCategory)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Category already exists"})
		return
	}

	category.Name = input.Name
	result, err := categoriesCollection.InsertOne(context.TODO(), category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create categories"})
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

func (app *application) getCategories(c *gin.Context) {
	categoriesCollection := app.Collection(data.CollectionCategory)
	var categories []bson.M

	cursor, err := categoriesCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching categories"})
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var category bson.M
		if err := cursor.Decode(&category); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding category"})
			return
		}
		categories = append(categories, category)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (app *application) updateCategory(c *gin.Context) {
	categoryID := c.Param("id")

	var input struct {
		Name string `json:"name"`
	}

	objectID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	categoriesCollection := app.Collection(data.CollectionCategory)

	if err := c.ShouldBindJSON((&input)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var existingCategory data.Category
	err = categoriesCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&existingCategory)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{"name": input.Name},
	}
	result, err := categoriesCollection.UpdateByID(context.TODO(), objectID, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating category"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No changes made"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}
