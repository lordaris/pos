package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (app *application) createProduct(c *gin.Context) {
	var input struct {
		Name        string  `json:"name"`
		Brand       string  `json:"brand"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		MinStock    int     `json:"min_stock"`
		Barcode     int     `json:"barcode"`
		CategoryID  string  `json:"category_id"`
	}

	// Bind the JSON body from the request to the `input` struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a product with the same barcode exists
	productsCollection := app.Collection(data.CollectionProduct)
	var existingBarcode data.Product
	err := productsCollection.FindOne(context.TODO(), bson.M{"barcode": input.Barcode}).Decode(&existingBarcode)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Barcode already exists", "product": existingBarcode.Name})
		return
	}

	// Check if it's a valid CategoryID
	categoryObjectID, err := primitive.ObjectIDFromHex(input.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Check if the category exists in database
	categoryCollection := app.Collection(data.CollectionCategory)
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
		Name:        input.Name,
		Brand:       input.Brand,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		MinStock:    input.MinStock,
		Barcode:     input.Barcode,
		CategoryID:  categoryObjectID,
	}

	// Insert the new product document into the product collection
	result, err := productsCollection.InsertOne(context.TODO(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

func (app *application) getProduct(c *gin.Context) {
	barcodeStr := c.Param("barcode")
	barcode, err := strconv.Atoi(barcodeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid barcode"})
	}

	productsCollection := app.Collection(data.CollectionProduct)
	filter := bson.D{{"barcode", barcode}}

	var existingProduct data.Product
	err = productsCollection.FindOne(context.TODO(), filter).Decode(&existingProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": existingProduct})
}

func (app *application) deleteProduct(c *gin.Context) {
	barcodeStr := c.Param("barcode")
	barcode, err := strconv.Atoi(barcodeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// Find product by Barcode
	productsCollection := app.Collection(data.CollectionProduct)
	filter := bson.D{{"barcode", barcode}}
	result, err := productsCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (app *application) updateProduct(c *gin.Context) {
	barcodeStr := c.Param("barcode")
	barcode, err := strconv.Atoi(barcodeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	var input struct {
		Name        *string  `json:"name"`
		Brand       *string  `json:"brand"`
		Description *string  `json:"description"`
		Price       *float64 `json:"price"`
		Stock       *int     `json:"stock"`
		MinStock    *int     `json:"min_stock"`
		Barcode     *int     `json:"barcode"`
		CategoryID  *string  `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productsCollection := app.Collection(data.CollectionProduct)

	var existingBarcode data.Product
	err = productsCollection.FindOne(context.TODO(), bson.M{"barcode": input.Barcode}).Decode(&existingBarcode)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Barcode already exists", "product": existingBarcode.Name})
		return
	}

	var existingProduct data.Product
	err = productsCollection.FindOne(context.TODO(), bson.M{"barcode": barcode}).Decode(&existingProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{}
	if input.Name != nil {
		update["name"] = *input.Name
	}
	if input.Brand != nil {
		update["brand"] = *input.Brand
	}
	if input.Description != nil {
		update["description"] = *input.Description
	}
	if input.Price != nil {
		update["price"] = *input.Price
	}
	if input.Stock != nil {
		update["stock"] = *input.Stock
	}
	if input.MinStock != nil {
		update["min_stock"] = *input.MinStock
	}
	if input.Barcode != nil {
		update["barcode"] = *input.Barcode
	}
	if input.CategoryID != nil {
		update["category_id"] = *input.CategoryID
	}

	if len(update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	var updatedProduct data.Product
	err = productsCollection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"barcode": barcode},
		bson.M{"$set": update},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": updatedProduct})
}
