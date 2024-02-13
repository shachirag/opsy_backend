package logEntry

type EditMentalHealthReqDto struct {
	Notes string   `json:"notes" bson:"notes"`
	When  string   `json:"when" bson:"when"`
	Ways  []string `json:"ways" bson:"ways"`
	Feel  string   `json:"feel" bson:"feel"`
}

type EditPhysicalHealthReqDto struct {
	Notes     string   `json:"notes" bson:"notes"`
	When      string   `json:"when" bson:"when"`
	PainLevel int32    `json:"painLevel" bson:"painlevel"`
	Ways      []string `json:"ways" bson:"ways"`
}

type EditPhysicalHealthResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type EditMentalHealthResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
