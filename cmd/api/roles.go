package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lordaris/pos-api/cmd/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (app *application) defaultRoles(c *gin.Context) {
	adminRole := data.Role{
		Name:        "admin",
		Permissions: []string{"all_permissions"},
	}

	managerRole := data.Role{
		Name: "manager",
		Permissions: []string{
			"manage_users", "manage_inventory", "data_visualization", "management_reports",
			"purchase_request_authorization", "inventory_transfer_authorization",
		},
	}

	cashierRole := data.Role{
		Name:        "cashier",
		Permissions: []string{"pos"},
	}

	roles := []data.Role{adminRole, managerRole, cashierRole}
	collection := app.config.db.mongoClient.Database("pos").Collection("roles")

	for _, role := range roles {
		filter := bson.M{"name": role.Name}
		update := bson.M{
			"$set": bson.M{
				"permissions": role.Permissions,
			},
		}

		opts := options.Update().SetUpsert(true)
		_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				// No document found, so let's insert a new role
				_, err := collection.InsertOne(context.TODO(), role)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role: " + role.Name})
					return
				}
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role: " + role.Name})
				return
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Roles created/updated successfully"})
}
