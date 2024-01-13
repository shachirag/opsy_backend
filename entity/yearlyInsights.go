package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type YearlyInsightsEntity struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	MentalHealth   []MentalHealth     `json:"mentalHealth" bson:"mentalHealth"`
	PhysicalHealth []PhysicalHealth   `json:"physicalHealth" bson:"physicalHealth"`
	Month          int32              `json:"month" bson:"month"`
	Year           int32              `json:"year" bson:"year"`
}

type MentalHealth struct {
	AvgFeel              int32 `json:"avgFeel" bson:"avgFeel"`
	TotalMentalHealthLog int64 `json:"totalMentalHealthLog" bson:"totalMentalHealthLog"`
}

type PhysicalHealth struct {
	AvgPain              int32 `json:"avgPain" bson:"avgPain"`
	TotalMentalHealthLog int64 `json:"totalMentalHealthLog" bson:"totalMentalHealthLog"`
}
