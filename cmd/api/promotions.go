package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *application) productPromotion(c *gin.Context) {
	var input struct {
		Barcode            int       `json:"barcode"`
		Type               string    `json:"type"`
		DiscountPercentage *int      `json:"discount_percentage"`
		DiscountPrice      *float32  `json:"discount_price"`
		BuyQuantity        *int      `json:"buy_quantity"`
		GetQuantity        *int      `json:"get_quantity"`
		StartDate          time.Time `json:"start_date"`
		EndDate            time.Time `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productsCollection := app.Collection(data.CollectionProduct)
	filter := bson.D{{"barcode", input.Barcode}}
	var existingProduct data.Product
	err := productsCollection.FindOne(context.TODO(), filter).Decode(&existingProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Validate startdate
	if input.StartDate.After(input.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "StartDate should be less than EndDate"})
		return
	}

	switch input.Type {
	case "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "type cannot be empty, it should contain DiscountPercentage | DiscountPrice | BuyGet"})
		return
	case "DiscountPercentage":
		if input.DiscountPercentage == nil || *input.DiscountPercentage <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "DiscountPercentage should be provided and greater than zero"})
			return
		} else if *input.DiscountPercentage >= 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "DiscountPercentage should be less than 100"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"DiscountPercentage": input.DiscountPercentage})
	case "DiscountPrice":
		if input.DiscountPrice == nil || *input.DiscountPrice == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "DiscountPrice should be provided and greater than 0"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"DiscountPrice": input.DiscountPercentage})
	case "BuyGet":
		if input.BuyQuantity == nil || *input.BuyQuantity <= 0 || input.GetQuantity == nil || *input.GetQuantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "BuyQuantity and GetQuantity should be provided and greater than 0"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid promotion type"})
		return
	}

	newPromotion := &data.Promotion{
		ID:        primitive.NewObjectID(),
		Type:      input.Type,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
	}

	if input.DiscountPercentage != nil {
		newPromotion.DiscountPercentage = *input.DiscountPercentage
	}

	if input.DiscountPrice != nil {
		newPromotion.DiscountPrice = *input.DiscountPrice
	}

	if input.BuyQuantity != nil {
		newPromotion.BuyQuantity = *input.BuyQuantity
	}

	if input.GetQuantity != nil {
		newPromotion.GetQuantity = *input.GetQuantity
	}

	update := bson.D{{"$set", bson.D{
		{"promotion", newPromotion},
	}}}
	_, err = productsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": newPromotion})
}

func (app *application) getPromotion(c *gin.Context) {
	barcodeStr := c.Param("barcode")
	barcode, err := strconv.Atoi(barcodeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusOK, gin.H{"promotion": existingProduct.Promotion, "product": existingProduct.Name, "barcode": existingProduct.Barcode})
}
