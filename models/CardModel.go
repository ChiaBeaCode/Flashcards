package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CardModel struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	// ID primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
	Title      string `json:"title,omitempty"`
	Definition string `json:"definition,omitempty"`
}
