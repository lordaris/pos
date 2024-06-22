package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Promotion struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	Type               string             `bson:"type"`
	DiscountPercentage float64            `bson:"discount_percentage,omitempty"`
	DiscountPrice      float64            `bson:"discount_price,omitempty"`
	BuyQuantity        int                `bson:"buy_quantity,omitempty"`
	GetQuantity        int                `bson:"get_quantity,omitempty"`
	StartDate          time.Time          `bson:"start_date"`
	EndDate            time.Time          `bson:"end_date"`
}
