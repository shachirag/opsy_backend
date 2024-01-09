package logEntry

// type InsightResDto struct {
// 	Status  bool         `json:"status"`
// 	Message string       `json:"message"`
// 	Data    []InsightRes `json:"data"`
// }

// type MentalHealthInsightResDto struct {
// 	Status  bool                     `json:"status"`
// 	Message string                   `json:"message"`
// 	Data    []MentalHealthInsightRes `json:"data"`
// }

type InsightsResDto struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    InsightsData `json:"data"`
}

type InsightsData struct {
	MentalHealth   []MentalHealthRes   `json:"mentalHealth"`
	PhysicalHealth []PhysicalHealthRes `json:"physicalHealth"`
}

type PhysicalHealthRes struct {
	Date         string   `json:"date" bson:"date"`
	Day          string   `json:"day" bson:"day"`
	AvgPainLevel *float64 `json:"avgPainLevel" bson:"avgPainLevel"`
}

type MentalHealthRes struct {
	Date    string  `json:"date" bson:"date"`
	Day     string  `json:"day" bson:"day"`
	AvgFeel float64 `json:"avgFeel" bson:"avgFeel"`
}
