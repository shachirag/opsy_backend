package userAuth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetUserResDto struct {
	Status  bool     `json:"status"`
	Message string   `json:"message"`
	User    UserData `json:"data"`
}

type UserData struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	Name          string             `json:"name" bson:"name"`
	Email         string             `json:"email" bson:"email"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
}