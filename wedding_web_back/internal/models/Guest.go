package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Guest struct {
	Id          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `json:"name"`
	Number      string             `json:"number"`
	Disposition string             `json:"disposition"`
}
