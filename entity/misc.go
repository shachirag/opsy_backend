package entity

type MiscEntity struct {
	Id      string   `json:"id" bson:"_id"`
	Popular []string `json:"popular" bson:"popular"`
	Other   []string `json:"other" bson:"other"`
}
