package userAuth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginReqDto struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type LoginResDto struct {
	Status  bool       `json:"status" bson:"status"`
	Message string     `json:"message" bson:"message"`
	Data    UserResDto `json:"data" bson:"data"`
	Token   string     `json:"token" bson:"token"`
}

type UserResDto struct {
	Id        primitive.ObjectID `json:"id" bson:"id"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
