package userAuth

type CatgoriesResDto struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Data    CategoriesRes `json:"data"`
}

type CategoriesRes struct {
	PhycicalHealth PhycicalHealth `json:"physicalHealth"`
	MentalHealth   MentalHealth   `json:"mentalHealth"`
}

type PhycicalHealth struct {
	Popular []string `json:"popular"`
	Other   []string `json:"other"`
}

type MentalHealth struct {
	Popular []string `json:"popular"`
	Other   []string `json:"other"`
}

type Misc struct {
	ID      string   `bson:"_id"`
	Popular []string `bson:"popular"`
	Other   []string `bson:"other"`
}
