package userAuth

type SignupReqDto struct {
	Email string `json:"email" bson:"email"`
}

type SignupResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type VerifyOtpSignupReqDto struct {
	Name       string `json:"name" bson:"name"`
	Email      string `json:"email" bson:"email"`
	Password   string `json:"password" bson:"password"`
	EnteredOTP string `json:"otp" bson:"otp"`
}

type VerifyOtpSignupResDto struct {
	Status  bool       `json:"status" bson:"status"`
	Message string     `json:"message" bson:"message"`
	Data    UserResDto `json:"data" bson:"data"`
	Token   string     `json:"token" bson:"token"`
}
