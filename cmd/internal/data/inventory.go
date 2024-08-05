package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionInventory = "inventory_movement"
)

// Type can be an adjustment, entry of new products, product transfer, etc...
type InventoryMovement struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Barcode      int                `bson:"barcode"`
	MovementType string             `bson:"movement_type"`
	Quantity     int                `bson:"quantity"`
	MovementDate time.Time          `bson:"movement_date"`
	UserID       primitive.ObjectID `bson:"user_id,omitempty"`
	Reason       string             `bson:"reason,omitempty"`
}
