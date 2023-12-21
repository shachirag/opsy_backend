package userAuth

type UpdateUserReqDto struct {
	Name              string   `json:"name" form:"name"`
	Email              string   `json:"email" form:"email"`
	OldProfileImage   string   `json:"oldUserImage" form:"oldUserImage"`
	OldQualifications []string `json:"oldQualifications" form:"oldQualifications"`
}

type UpdateUserResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}