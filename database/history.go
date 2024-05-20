package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type History struct {
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	Name               string             `bson:"name" json:"name"`
	Service            string             `bson:"service" json:"service"`
	Type               string             `bson:"type" json:"type"`
	AIModel            string             `bson:"aiModel" json:"aiModel"`
	AIConfiguration    string             `bson:"aiConfiguration" json:"aiConfiguration"`
	Prompt             string             `bson:"prompt" json:"prompt"`
	Description        string             `bson:"description" json:"description"`
	Style              string             `bson:"style" json:"style"`
	ColorScheme        string             `bson:"colorScheme" json:"colorScheme"`
	Text               string             `bson:"text" json:"text"`
	TextStyle          string             `bson:"textStyle" json:"textStyle"`
	Layout             string             `bson:"layout" json:"layout"`
	Theme              string             `bson:"theme" json:"theme"`
	AdditionalElements string             `bson:"additionalElements" json:"additionalElements"`
	Product            string             `bson:"product" json:"product"`
	CreatedAt          time.Time          `bson:"createdAt" json:"createdAt"`
}
