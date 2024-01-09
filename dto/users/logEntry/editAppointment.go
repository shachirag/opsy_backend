package logEntry

type GetAppointmentResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type AppointmentReqDto struct {
	Notes       string `json:"notes" bson:"notes"`
	When        string `json:"when" bson:"when"`
	WhatItIsFor string `json:"whatItIsFor" bson:"whatItIsFor"`
	Alert       string `json:"alert" bson:"alert"`
}
