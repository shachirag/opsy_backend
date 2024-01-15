package logEntry

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetAppointmentResDto struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    AppointmentRes `json:"data"`
}

type AppointmentRes struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Notes       string             `json:"notes" bson:"notes"`
	When        time.Time          `json:"when" bson:"when"`
	WhatItIsFor string             `json:"whatItIsFor" bson:"whatItIsFor"`
	Alert       string             `json:"alert" bson:"alert"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type FutureAppointmentDto struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    []FutureAppointmentRes `json:"data"`
}

type FutureAppointmentRes struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Notes       string             `json:"notes" bson:"notes"`
	When        time.Time          `json:"when" bson:"when"`
	WhatItIsFor string             `json:"whatItIsFor" bson:"whatItIsFor"`
	Alert       string             `json:"alert" bson:"alert"`
	NumberCount int64              `json:"numberCount" bson:"numberCount"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type AppointmentReqDto struct {
	Notes       string `json:"notes" bson:"notes"`
	When        string `json:"when" bson:"when"`
	WhatItIsFor string `json:"whatItIsFor" bson:"whatItIsFor"`
	Alert       string `json:"alert" bson:"alert"`
}
