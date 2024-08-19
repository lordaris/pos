package main

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
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

		// Store the invoice items with its information
		invoiceItems := make([]data.InvoiceItem, 0, len(input.Items))

		for _, item := range input.Items {
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

			invoiceItem := data.InvoiceItem{
				ProductName: product.Name,
				Barcode:     product.Barcode,
				Quantity:    item.Quantity,
				Price:       product.Price,
			}

			if product.Promotion.StartDate.Before(time.Now()) && product.Promotion.EndDate.After(time.Now()) {
				invoiceItem.PromotionType = product.Promotion.Type
				switch product.Promotion.Type {
				case "DiscountPercentage":
					invoiceItem.PriceWithDiscount = math.Round((product.Price-(product.Price*float64(product.Promotion.DiscountPercentage)/100))*100) / 100
					invoiceItem.DiscountPercentage = product.Promotion.DiscountPercentage
					invoiceItem.TotalAmount = invoiceItem.PriceWithDiscount * float64(item.Quantity)
					totalAmount += float64(invoiceItem.TotalAmount)

				case "DiscountPrice":
					invoiceItem.DiscountPrice = float64(product.Promotion.DiscountPrice)
					invoiceItem.TotalAmount = invoiceItem.DiscountPrice * float64(item.Quantity)
					totalAmount += invoiceItem.TotalAmount

				case "BuyGet":
					productModule := item.Quantity % product.Promotion.GetQuantity
					itemGetQuotient := int(item.Quantity / product.Promotion.GetQuantity)
					invoiceItem.PaidQuantity = (itemGetQuotient * product.Promotion.BuyQuantity) + productModule
					invoiceItem.FreeQuantity = item.Quantity - invoiceItem.PaidQuantity
					invoiceItem.TotalAmount = product.Price * float64(invoiceItem.PaidQuantity)

					totalAmount += invoiceItem.TotalAmount
				}
			} else {
				invoiceItem.TotalAmount = product.Price * float64(item.Quantity)
				totalAmount += invoiceItem.TotalAmount

			}

			invoiceItems = append(invoiceItems, invoiceItem)

		}

		totalAmount = math.Round(totalAmount*100) / 100

		// Create the invoice object
		invoice := &data.Invoice{
			ID:           primitive.NewObjectID(),
			TicketNumber: ticketNumber,
			TotalAmount:  totalAmount,
			SaleDate:     time.Now(),
			Items:        invoiceItems,
			UserID:       user.ID,
			Username:     user.Name,
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

func (app *application) getInvoice(c *gin.Context) {
	ticketNumberStr := c.Param("ticketnumber")
	ticketNumber, err := strconv.Atoi(ticketNumberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket number"})
	}

	invoicesCollection := app.Collection(data.CollectionInvoice)
	filter := bson.D{{"ticket_number", ticketNumber}}
	var existingInvoice data.Invoice
	err = invoicesCollection.FindOne(context.TODO(), filter).Decode(&existingInvoice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Ticket": existingInvoice})
}
