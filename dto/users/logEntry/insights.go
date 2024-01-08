package logEntry

type InsightResDto struct {
	Status  bool         `json:"status"`
	Message string       `json:"message"`
	Data    []InsightRes `json:"data"`
}

type InsightRes struct {
	Date         string   `json:"date" bson:"date"`
	Day          string   `json:"day" bson:"day"`
	AvgPainLevel *float64 `json:"avgPainLevel" bson:"avgPainLevel"`
}


type MentalHealthInsightResDto struct {
    Status  bool                    `json:"status"`
    Message string                  `json:"message"`
    Data    []MentalHealthInsightRes `json:"data"`
}
 
type MentalHealthInsightRes struct {
    Date    string  `json:"date" bson:"date"`
    Day     string  `json:"day" bson:"day"`
    AvgFeel float64 `json:"avgFeel" bson:"avgFeel"`
}



