package data

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionCategory = "categories"
)

type Category struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}
