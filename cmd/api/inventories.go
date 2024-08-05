package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Constants for movement types
const (
	TypeIn  = "IN"
	TypeOut = "OUT"
)

func (app *application) createInventoryMovement(c *gin.Context) {
	user := app.contextGetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not authenticated"})
		return
	}

	var input struct {
		Barcode      int    `json:"barcode"`
		MovementType string `json:"movement_type"`
		Quantity     int    `json:"quantity"`
		Reason       string `json:"reason,omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate type of movement
	if input.MovementType != TypeIn && input.MovementType != TypeOut {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movement type"})
		return
	}

	// Start session for transaction
	session, err := app.config.db.mongoClient.StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
		return
	}

	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(sc mongo.SessionContext) (interface{}, error) {
		productsCollection := app.Collection(data.CollectionProduct)
		filter := bson.M{"barcode": input.Barcode}
		update := bson.M{}
		if input.MovementType == TypeIn {
			update["$inc"] = bson.M{"stock": input.Quantity}
		} else {
			update["$inc"] = bson.M{"stock": -input.Quantity}
		}

		result, err := productsCollection.UpdateOne(sc, filter, update)
		if err != nil {
			return nil, err
		}
		if result.MatchedCount == 0 {
			return nil, mongo.ErrNoDocuments
		}

		// Initialize inventory_movements collection
		movementsCollection := app.Collection(data.CollectionInventory)
		movement := data.InventoryMovement{
			Barcode:      input.Barcode,
			MovementType: input.MovementType,
			Quantity:     input.Quantity,
			MovementDate: time.Now(),
			UserID:       user.ID,
			Reason:       input.Reason,
		}

		_, err = movementsCollection.InsertOne(sc, movement)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, "Inventory movement created successfully")
}
