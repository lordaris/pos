package data

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Brand       string             `bson:"brand"`
	Description string             `bson:"description"`
	Price       float64            `bson:"price"`
	Stock       int                `bson:"stock"`
	MinStock    int                `bson:"min_stock"`
	Barcode     string             `bson:"barcode"`
	PLU         int                `bson:"plu"`
	CategoryID  primitive.ObjectID `bson:"category_id"`
	Promotion   Promotion          `bson:"promotion,omitempty"`
}
