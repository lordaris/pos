package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Barcode string             `bson:"barcode"`
	Name    string             `bson:"name"`
	Email   string             `bson:"email"`
	Phone   string             `bson:"phone"`
	Points  int                `bson:"points"`
}
