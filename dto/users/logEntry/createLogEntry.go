package logEntry

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetLogEntryResDto struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    LogEntryRes `json:"data"`
}

type LogEntryRes struct {
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

type LogEntryReqDto struct {
	Type        string             `json:"type" bson:"type"`
	Feel        string             `json:"feel" bson:"feel"`
	Notes       string             `json:"notes" bson:"notes"`
	Ways        []string           `json:"ways" bson:"ways"`
	DateTime    string             `json:"dateTime" bson:"dateTime"`
	PainLevel   int32              `json:"painLevel" bson:"painlevel"`
	WhatItIsFor string             `json:"whatItIsFor" bson:"whatItIsFor"`
	Alert       string             `json:"alert" bson:"alert"`
	UserId      primitive.ObjectID `json:"userId" bson:"userId"`
}
