package Entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type MongoBenchEntity struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Name    string             `bson:"Name"`
	Counter int                `bson:"Counter"`
}
