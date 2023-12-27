package userAuth

type CatgoriesResDto struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Data    CategoriesRes `json:"data"`
}

type CategoriesRes struct {
	PhysicalHealth PhysicalHealth`json:"physicalHealth"`
	MentalHealth   MentalHealth   `json:"mentalHealth"`
}

type PhysicalHealth struct {
	Popular []string `json:"popular"`
	Other   []string `json:"other"`
}

type MentalHealth struct {
	Popular []string `json:"popular"`
	Other   []string `json:"other"`
}