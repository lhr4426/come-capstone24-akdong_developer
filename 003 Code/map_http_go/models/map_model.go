package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Map struct {
	Map_ID primitive.ObjectID `json:id`
	Version int `json:version,omitempty`
	Chunk int `json:version`
}