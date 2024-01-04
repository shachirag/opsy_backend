package logEntry

type GetLogEntryResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type LogEntryReqDto struct {
	Type        string   `json:"type" bson:"type"`
	Feel        string   `json:"feel" bson:"feel"`
	Notes       string   `json:"notes" bson:"notes"`
	Ways        []string `json:"ways" bson:"ways"`
	When        string   `json:"when" bson:"when"`
	PainLevel   int32    `json:"painLevel" bson:"painlevel"`
	WhatItIsFor string   `json:"whatItIsFor" bson:"whatItIsFor"`
	Alert       string   `json:"alert" bson:"alert"`
}
