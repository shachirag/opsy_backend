 package logEntry

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CatgoriesResDto struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Data    []CategoriesRes `json:"data"`
}

type CategoriesRes struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Type        string             `json:"type" bson:"type"`
	Feel        string             `json:"feel" bson:"feel"`
	When        time.Time          `json:"when" bson:"when"`
	PainLevel   int32              `json:"painLevel" bson:"painLevel"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
