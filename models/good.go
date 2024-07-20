package models

// Goods represents an item in the contract.
type Goods struct {
	Item     string `bson:"item"`
	Quantity int    `bson:"quantity"`
	Unit     string `bson:"unit"`
	Price    int    `bson:"price"`
}
