package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryMovement struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ProductID    primitive.ObjectID `bson:"product_id"`
	MovementType string             `bson:"movement_type"`
	Quantity     int                `bson:"quantity"`
	MovementDate time.Time          `bson:"movement_date"`
	UserID       primitive.ObjectID `bson:"user_id,omitempty"`
	Reason       string             `bson:"reason,omitempty"`
}

type Customer struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Barcode string             `bson:"barcode"`
	Name    string             `bson:"name"`
	Email   string             `bson:"email"`
	Phone   string             `bson:"phone"`
	Points  int                `bson:"points"`
}
