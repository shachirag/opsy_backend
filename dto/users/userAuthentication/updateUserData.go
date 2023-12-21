package userAuth

type UpdateUserReqDto struct {
	Name string `json:"name" form:"name"`
}

type UpdateUserResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
