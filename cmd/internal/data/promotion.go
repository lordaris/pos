package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Promotion struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	Type               string             `bson:"type"`
	DiscountPercentage int                `bson:"discount_percentage"`
	DiscountPrice      float64            `bson:"discount_price"`
	BuyQuantity        int                `bson:"buy_quantity"`
	GetQuantity        int                `bson:"get_quantity"`
	StartDate          time.Time          `bson:"start_date"`
	EndDate            time.Time          `bson:"end_date"`
}
