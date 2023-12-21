package userAuth

type ChangeUserPasswordReqDto struct {
	CurrentPassword string `json:"currentPassword" bson:"currentPassword"`
	NewPassword     string `json:"newPassword" bson:"newPassword"`
	ConfirmPassword string `json:"confirmPassword" bson:"confirmPassword"`
}

type ChangeUserPasswordResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
