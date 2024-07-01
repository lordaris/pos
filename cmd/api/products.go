package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: Implement the ability to store an array of offers within a product document.
// These offers should not have overlapping date ranges, allowing to schedule various promotions in advance.
// The point-of-sale (POS) system logic should account for this and only apply active offers during transactions.
// **Additionally, a mechanism should be implemented to automatically remove expired offers from the product document.**

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
	productsCollection := app.config.db.mongoClient.Database("pos").Collection("products")
	var existingBarcode data.Product
	err := productsCollection.FindOne(context.TODO(), bson.M{"barcode": input.Barcode}).Decode(&existingBarcode)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Barcode already exists"})
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
		CategoryID:  categoryObjectID,
	}

	// Insert the new user document into the user collection
	result, err := productsCollection.InsertOne(context.TODO(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

func (app *application) productPromotion(c *gin.Context) {
	// TODO: Check how to handle the creation of the start date and end date
	type promotion struct {
		Type               string     `json:"type"`
		DiscountPercentage *int       `json:"discount_percentage"`
		DiscountPrice      *float64   `json:"discount_price"`
		BuyQuantity        *int       `json:"buy_quantity"`
		GetQuantity        *int       `json:"get_quantity"`
		StartDate          *time.Time `json:"start_date"`
		EndDate            *time.Time `json:"end_date"`
	}

	// Get the product by PLU or barcode
	var input struct {
		Barcode   int `json:"barcode"`
		Promotion promotion
	}

	// Bind the JSON body from the request to the `input` struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productCollection := app.config.db.mongoClient.Database("pos").Collection("products")
	filter := bson.D{{"barcode", input.Barcode}}
	var existingProduct data.Product
	err := productCollection.FindOne(context.TODO(), filter).Decode(&existingProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch input.Promotion.Type {
	case "DiscountPercentage":
		if input.Promotion.DiscountPercentage == nil || *input.Promotion.DiscountPercentage <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "DiscountPercentage should be provided and greater than zero"})
			return
		} else if *input.Promotion.DiscountPercentage >= 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "DiscountPercentage should be less than 100"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"DiscountPercentage": input.Promotion.DiscountPercentage})
	case "DiscountPrice":
		if input.Promotion.DiscountPrice == nil || *input.Promotion.DiscountPrice == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "DiscountPrice should be provided and greater than 0"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"DiscountPrice": input.Promotion.DiscountPercentage})
	case "BuyGet":
		if input.Promotion.BuyQuantity == nil || *input.Promotion.BuyQuantity <= 0 || input.Promotion.GetQuantity == nil || *input.Promotion.GetQuantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "BuyQuantity and GetQuantity should be provided and greater than 0"})
			return
		}

	default:
		fmt.Println("You should select one of the next options: DiscountPercentage | DiscountPrice | BuyGet")
	}

	update := bson.D{{"$set", bson.D{
		{"promotion", input.Promotion},
	}}}
	_, err = productCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product with promotion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion applied successfully"})
}
