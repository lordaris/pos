package data

import (
	"reflect"

	"github.com/lordaris/pos-api/cmd/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionProduct = "products"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Brand       string             `bson:"brand"`
	Description string             `bson:"description"`
	Price       float64            `bson:"price"`
	Stock       int                `bson:"stock"`
	MinStock    int                `bson:"min_stock"`
	Barcode     int                `bson:"barcode"`
	CategoryID  primitive.ObjectID `bson:"category_id"`
	Promotion   Promotion          `bson:"promotion"`
}

func validatePoduct(v *validator.Validator, product *Product) {
	v.Check(product.Name != "", "name", "must be provided")
	v.Check(product.Brand != "", "brand", "must be provided")
	v.Check(product.Description != "", "description", "must be provided")
	v.Check(reflect.TypeOf(product.Price).Kind() == reflect.Float64, "price", "must be a float64")
}
