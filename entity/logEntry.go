package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogEntryEntity struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Type        string             `json:"type" bson:"type"`
	Feel        string             `json:"feel" bson:"feel"`
	IsDeleted   bool               `json:"isDeleted" bson:"isDeleted"`
	Notes       string             `json:"notes" bson:"notes"`
	Ways        []string           `json:"ways" bson:"ways"`
	When        time.Time          `json:"when" bson:"when"`
	PainLevel   int32              `json:"painLevel" bson:"painlevel"`
	WhatItIsFor string             `json:"whatItIsFor" bson:"whatItIsFor"`
	Alert       string             `json:"alert" bson:"alert"`
	UserId          primitive.ObjectID `json:"userId" bson:"userId"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
