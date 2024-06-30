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

func (app *application) createProduct(c *gin.Context) {
	var input struct {
		Name        string  `json:"name"`
		Brand       string  `json:"brand"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		MinStock    int     `json:"min_stock"`
		Barcode     string  `json:"barcode"`
		PLU         int     `json:"plu"`
		CategoryID  string  `json:"category_id"`
	}

	// Bind the JSON body from the request to the `input` struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a product with the same name exists
	productsCollection := app.config.db.mongoClient.Database("pos").Collection("products")
	var existingProduct data.Product
	err := productsCollection.FindOne(context.TODO(), bson.M{"name": input.Name}).Decode(&existingProduct)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product already exists"})
		return
	}

	// Check if it's a valid CategoryID
	categoryObjectID, err := primitive.ObjectIDFromHex(input.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Check if the category exists in database
	categoryCollection := app.config.db.mongoClient.Database("pos").Collection("categories")
	var existingCategory data.Category
	err = categoryCollection.FindOne(context.TODO(), bson.M{"_id": categoryObjectID}).Decode(&existingCategory)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	product := &data.Product{
		ID:          primitive.NewObjectID(),
		Name:        input.Name,
		Brand:       input.Brand,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		MinStock:    input.MinStock,
		Barcode:     input.Barcode,
		PLU:         input.PLU,
		CategoryID:  categoryObjectID,
	}

	// Insert the new user document into the user collection
	result, err := productsCollection.InsertOne(context.TODO(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}
