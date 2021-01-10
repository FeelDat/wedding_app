package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Login struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	Login    string             `json:"login"`
	Password string             `json:"password"`
}
