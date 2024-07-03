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
)

func (app *application) createInvoice(c *gin.Context) {
	var input struct {
		TotalAmount float64 `json:"total_amount"`
		//	UserID      primitive.ObjectID    `json:"user_id"`
		// CustomerID  primitive.ObjectID    `json:"customer_id,omitempty"`
		SaleDate time.Time `json:"sale_date"`
		//	ChangeGiven float64               `json:"change_given"`
		//	Discount    float64               `json:"discount,omitempty"`
		Items []struct {
			ProductID primitive.ObjectID `json:"product_id"`
			Quantity  int                `json:"quantity"`
			Price     float64            `json:"price,omitempty"`
		} `json:"items"`
		// 	Payments    []data.InvoicePayment `json:"payments"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoiceCollection := app.config.db.mongoClient.Database("pos").Collection("invoices")
	productsCollection := app.config.db.mongoClient.Database("pos").Collection("products")

	// Calculate total amount considering promotions.
	var totalAmount float64
	for i, item := range input.Items {
		var product data.Product
		filter := bson.M{"_id": item.ProductID}
		err := productsCollection.FindOne(context.TODO(), filter).Decode(&product)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
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
					// Calculate the remainder when dividing the quantity of items by the get quantity
					// specified in the promotion. It determine how many items are left out of the promotion.
					productModule := item.Quantity % promotion.GetQuantity
					// represent how many full sets of get quantities fit into the total quantity of items.
					itemGetQuotient := int(item.Quantity / promotion.GetQuantity)
					// Compute the total number of products that need to be paid for (itemGetQuotient * promotion.BuyQuantity) and any
					// additional products that do not meet the promotion criteria.
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
		ID:          primitive.NewObjectID(),
		TotalAmount: totalAmount,
		SaleDate:    time.Now(),
		Items:       make([]data.InvoiceItem, len(input.Items)),
	}

	// Assign each item to the invoice with its calculated price
	for i, item := range input.Items {
		invoice.Items[i] = data.InvoiceItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
	}
	result, err := invoiceCollection.InsertOne(context.TODO(), invoice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{"invoice": result})
}
