package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OtpEntity struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Email     string             `json:"email" bson:"email"`
	Otp       string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}
