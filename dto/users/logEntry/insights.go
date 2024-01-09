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
	Status  bool         `json:"status"`
	Message string       `json:"message"`
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

type MonthlyInsightsResDto struct {
	Status  bool                `json:"status"`
	Message string              `json:"message"`
	Data    MonthlyInsightsData `json:"data"`
}

type MonthlyInsightsData struct {
	MentalHealth   []MonthlyMentalHealthRes   `json:"mentalHealth"`
	PhysicalHealth []MonthlyPhysicalHealthRes `json:"physicalHealth"`
}

type MonthlyPhysicalHealthRes struct {
	Date         string   `json:"date" bson:"date"`
	AvgPainLevel *float64 `json:"avgPainLevel,omitempty" bson:"avgPainLevel,omitempty"`
}

type MonthlyMentalHealthRes struct {
	Date    string  `json:"date" bson:"date"`
	AvgFeel float64 `json:"avgFeel" bson:"avgFeel"`
}