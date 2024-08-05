package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (app *application) createInvoice(c *gin.Context) {
	user := app.contextGetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input struct {
		SaleDate time.Time `json:"sale_date"`
		Items    []struct {
			ProductName string  `json:"product_name,omitempty"`
			Barcode     int     `json:"barcode"`
			Quantity    int     `json:"quantity"`
			Price       float64 `json:"price,omitempty"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoiceCollection := app.Collection(data.CollectionInvoice)
	productsCollection := app.Collection(data.CollectionProduct)

	session, err := app.config.db.mongoClient.StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer session.EndSession(context.TODO())

	// Add numbering system, so a ticket number is returned after the invoice posting.
	// This number can be used to search for an specific ticket instead of searching for an ID, and also
	// adding that number to a physical ticket.
	// Using transactions (with mongo.WithSession) allows to execute multiple operations as a single logical unit of work,
	// and allow data consistency an integrity.
	var invoiceInfo data.Invoice
	_, err = session.WithTransaction(context.Background(), func(sc mongo.SessionContext) (interface{}, error) {
		// Find the last ticket number
		opts := options.FindOne().SetSort(bson.D{{"ticket_number", -1}})
		var lastInvoice data.Invoice
		err := invoiceCollection.FindOne(sc, bson.D{}, opts).Decode(&lastInvoice)
		if err != nil && err != mongo.ErrNoDocuments {
			return nil, err
		}

		ticketNumber := 1
		if err == nil {
			ticketNumber = lastInvoice.TicketNumber + 1
		}

		// Calculate total amount considering promotions
		var totalAmount float64
		for i, item := range input.Items {
			var product data.Product
			filter := bson.M{"barcode": item.Barcode}
			err := productsCollection.FindOne(sc, filter).Decode(&product)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
					return nil, fmt.Errorf("product not found")
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return nil, err // Return nil to abort the session without an error
			}

			// Update product stock
			updateStock := bson.M{"$inc": bson.M{"stock": -item.Quantity}}
			_, err = productsCollection.UpdateOne(sc, filter, updateStock)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
				return nil, fmt.Errorf("failed to update product stock")
			}

			price := product.Price
			if product.Promotion.StartDate.Before(time.Now()) && product.Promotion.EndDate.After(time.Now()) {
				switch product.Promotion.Type {
				case "DiscountPercentage":
					price = price - (price * float64(product.Promotion.DiscountPercentage) / 100)
					totalAmount = price * float64(item.Quantity)

				case "DiscountPrice":
					price = float64(product.Promotion.DiscountPrice)
					totalAmount = price * float64(item.Quantity)

				case "BuyGet":

					productModule := item.Quantity % product.Promotion.GetQuantity
					itemGetQuotient := int(item.Quantity / product.Promotion.GetQuantity)
					paidProducts := (itemGetQuotient * product.Promotion.BuyQuantity) + productModule
					totalAmount += product.Price * float64(paidProducts)

				}
			} else {
				totalAmount += price * float64(item.Quantity)
			}

			input.Items[i].Price = product.Price
			input.Items[i].ProductName = product.Name
		}

		totalAmount = math.Round(totalAmount*100) / 100
		// Create the invoice object
		invoice := &data.Invoice{
			ID:           primitive.NewObjectID(),
			TicketNumber: ticketNumber,
			TotalAmount:  totalAmount,
			SaleDate:     time.Now(),
			Items:        make([]data.InvoiceItem, len(input.Items)),
			UserID:       user.ID,
			Username:     user.Name,
		}

		for i, item := range input.Items {
			invoice.Items[i] = data.InvoiceItem{
				ProductName: item.ProductName,
				Barcode:     item.Barcode,
				Quantity:    item.Quantity,
				Price:       item.Price,
			}
		}

		invoiceInfo = *invoice

		// Insert the new invoice
		_, err = invoiceCollection.InsertOne(sc, invoice)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"invoice": invoiceInfo})
}
