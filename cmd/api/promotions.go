package main

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
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
		// TODO: Validate that BuyQuantity is less than GetQuantity
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

type promotionInfo struct {
	Type               *string    `json:"type"`
	DiscountPercentage *int       `json:"discount_percentage"`
	DiscountPrice      *float32   `json:"discount_price"`
	BuyQuantity        *int       `json:"buy_quantity"`
	GetQuantity        *int       `json:"get_quantity"`
	StartDate          *time.Time `json:"start_date"`
	EndDate            *time.Time `json:"end_date"`
}

// updatePromotion handles the update of a promotion for a specific product
func (app *application) updatePromotion(c *gin.Context) {
	barcodeStr := c.Param("barcode")
	barcode, err := strconv.Atoi(barcodeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid barcode"})
		return
	}

	// Define a variable with the promotionInfo structure for the promotion update fields

	var promotionUpdate promotionInfo

	// Bind the JSON from the request body to the structure
	if err := c.ShouldBindJSON(&promotionUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if at least one field is provided for update
	if !hasAnyUpdate(promotionUpdate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No promotion fields provided for update"})
		return
	}

	// Get the existing product document from the database
	productsCollection := app.Collection(data.CollectionProduct)
	var existingProduct data.Product
	err = productsCollection.FindOne(context.TODO(), bson.M{"barcode": barcode}).Decode(&existingProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		}
		return
	}

	// Check for existing promotions
	if existingProduct.Promotion.Type == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "The product doesn't have promotions"})
		return
	}

	// Validate and update the promotion
	if err := validateAndUpdatePromotion(&existingProduct.Promotion, promotionUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prepare the update for the database
	update := bson.M{"$set": bson.M{"promotion": existingProduct.Promotion}}

	// Update the document in the database
	_, err = productsCollection.UpdateOne(context.TODO(), bson.M{"barcode": barcode}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update promotion"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promotion updated successfully"})
}

// It checks if there's any value stored in the interface (if any field was provided for update)
func hasAnyUpdate(update interface{}) bool {
	val := reflect.ValueOf(update)
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).IsNil() {
			return true
		}
	}
	return false
}

// Validates and updates the promotion fields
func validateAndUpdatePromotion(existing *data.Promotion, update promotionInfo) error {
	if update.Type != nil {
		existing.Type = *update.Type
		switch *update.Type {
		case "BuyGet":
			if err := validateBuyGet(update.BuyQuantity, update.GetQuantity); err != nil {
				return err
			}
			existing.BuyQuantity = *update.BuyQuantity
			existing.GetQuantity = *update.GetQuantity
		case "DiscountPercentage":
			if err := validateDiscountPercentage(update.DiscountPercentage); err != nil {
				return err
			}
			existing.DiscountPercentage = *update.DiscountPercentage
		case "DiscountPrice":
			if err := validateDiscountPrice(update.DiscountPrice); err != nil {
				return err
			}
			existing.DiscountPrice = *update.DiscountPrice
		default:
			return fmt.Errorf("invalid promotion type")
		}
	}

	if update.StartDate != nil {
		existing.StartDate = *update.StartDate
	}
	if update.EndDate != nil {
		existing.EndDate = *update.EndDate
	}

	if existing.StartDate.After(existing.EndDate) {
		return fmt.Errorf("start date should be before end date")
	}

	return nil
}

func validateBuyGet(buy, get *int) error {
	if buy == nil || get == nil || *buy <= 0 || *get <= 0 {
		return fmt.Errorf("buy and get quantities should be provided and greater than 0")
	}
	if *buy >= *get {
		return fmt.Errorf("buy quantity should be less than get quantity")
	}
	return nil
}

func validateDiscountPercentage(discount *int) error {
	if discount == nil || *discount <= 0 || *discount >= 100 {
		return fmt.Errorf("discount percentage should be between 1 and 99")
	}
	return nil
}

func validateDiscountPrice(price *float32) error {
	if price == nil || *price <= 0 {
		return fmt.Errorf("discount price should be greater than zero")
	}
	return nil
}
