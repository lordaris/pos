package data

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	CollectionRole = "roles"
)

type Role struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}
