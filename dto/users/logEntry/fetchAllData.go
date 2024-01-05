package logEntry

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CatgoriesResDto struct {
	Status  bool            `json:"status"`
	Message string          `json:"message"`
	Data    []CategoriesRes `json:"data"`
}

type CategoriesRes struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Type        string             `json:"type" bson:"type"`
	Feel        string             `json:"feel" bson:"feel"`
	Notes       string             `json:"notes" bson:"notes"`
	Ways        []string           `json:"ways" bson:"ways"`
	When        time.Time          `json:"when" bson:"when"`
	PainLevel   int32              `json:"painLevel" bson:"painLevel"`
	WhatItIsFor string             `json:"whatItIsFor" bson:"whatItIsFor"`
	Alert       string             `json:"alert" bson:"alert"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// type When struct {
// 	Date string `json:"date" bson:"date"`
// 	Time string `json:"time" bson:"time"`
// }
