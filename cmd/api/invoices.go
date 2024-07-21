package main

import (
	"context"
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
	var input struct {
		TotalAmount float64   `json:"total_amount"`
		SaleDate    time.Time `json:"sale_date"`
		Items       []struct {
			ProductID primitive.ObjectID `json:"product_id"`
			Quantity  int                `json:"quantity"`
			Price     float64            `json:"price,omitempty"`
		} `json:"items"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoiceCollection := app.config.db.mongoClient.Database("pos").Collection("invoices")
	productsCollection := app.config.db.mongoClient.Database("pos").Collection("products")

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
	ticketNumber := 1
	err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		// Find the last ticket number
		opts := options.FindOne().SetSort(bson.D{{"ticket_number", -1}})
		var lastInvoice data.Invoice
		err := invoiceCollection.FindOne(sc, bson.D{}, opts).Decode(&lastInvoice)
		if err != nil && err != mongo.ErrNoDocuments {
			return err
		}

		if err == nil {
			ticketNumber = lastInvoice.TicketNumber + 1
		}

		// Calculate total amount considering promotions
		var totalAmount float64
		for i, item := range input.Items {
			var product data.Product
			filter := bson.M{"_id": item.ProductID}
			err := productsCollection.FindOne(sc, filter).Decode(&product)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
					return nil // Return nil to abort the session without an error
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return nil // Return nil to abort the session without an error
			}

			for _, promotion := range product.Promotions {
				if promotion.StartDate.Before(time.Now()) && promotion.EndDate.After(time.Now()) {
					switch promotion.Type {
					case "DiscountPercentage":
						itemPrice := product.Price * float64(item.Quantity)
						discount := itemPrice * float64(promotion.DiscountPercentage) / 100
						totalAmount += itemPrice - discount
					case "DiscountPrice":
						totalAmount += float64(item.Quantity) * float64(promotion.DiscountPrice)
					case "BuyGet":
						productModule := item.Quantity % promotion.GetQuantity
						itemGetQuotient := int(item.Quantity / promotion.GetQuantity)
						paidProducts := (itemGetQuotient * promotion.BuyQuantity) + productModule
						totalAmount += product.Price * float64(paidProducts)
					}
				} else {
					totalAmount += product.Price * float64(item.Quantity)
				}
			}

			if len(product.Promotions) == 0 {
				totalAmount += product.Price * float64(item.Quantity)
			}
			input.Items[i].Price = product.Price
		}
		// Create the invoice object
		invoice := &data.Invoice{
			ID:           primitive.NewObjectID(),
			TicketNumber: ticketNumber,
			TotalAmount:  totalAmount,
			SaleDate:     time.Now(),
			Items:        make([]data.InvoiceItem, len(input.Items)),
		}

		for i, item := range input.Items {
			invoice.Items[i] = data.InvoiceItem{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
			}
		}

		// Insert the new invoice
		_, err = invoiceCollection.InsertOne(sc, invoice)
		return err
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ticket_number": ticketNumber})
}
