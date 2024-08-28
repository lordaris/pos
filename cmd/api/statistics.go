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

func (app *application) getSales(c *gin.Context) {
	startDateString := c.Param("startdate")
	endDateString := c.Param("enddate")

	startDate, err := time.Parse(time.RFC3339, startDateString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse start date"})
		return
	}
	endDate, err := time.Parse(time.RFC3339, endDateString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse end date"})
		return
	}

	invoiceCollection := app.Collection(data.CollectionInvoice)

	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"sale_date", bson.D{{"$gte", startDate}, {"$lte", endDate}}}}}},
		bson.D{{"$group", bson.D{
			{"_id", nil},
			{"total_sales", bson.D{{"$sum", "$ticket_total"}}},
			{"total_invoices", bson.D{{"$sum", 1}}},
			{"average_sale", bson.D{{"$avg", "$ticket_total"}}},
			{"min_sale", bson.D{{"$min", "$ticket_total"}}},
			{"max_sale", bson.D{{"$max", "$ticket_total"}}},
		}}},
		bson.D{{"$project", bson.D{
			{"_id", 0},
			{"total_sales", bson.D{{"$round", bson.A{"$total_sales", 2}}}},
			{"total_invoices", 1},
			{"average_sale", bson.D{{"$round", bson.A{"$average_sale", 2}}}},
			{"min_sale", bson.D{{"$round", bson.A{"$min_sale", 2}}}},
			{"max_sale", bson.D{{"$round", bson.A{"$max_sale", 2}}}},
		}}},
	}

	cursor, err := invoiceCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute aggregation", "details": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode results", "details": err.Error()})
		return
	}

	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No sales found for the specified period"})
		return
	}

	c.JSON(http.StatusOK, results[0])
}
