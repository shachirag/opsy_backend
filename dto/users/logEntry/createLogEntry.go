package logEntry

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetLogEntryResDto struct {
	Status  bool               `json:"status"`
	Message string             `json:"message"`
	Id      primitive.ObjectID `json:"id"`
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
