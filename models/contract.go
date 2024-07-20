package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contract struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	CustomerID   primitive.ObjectID `bson:"customer_id"`
	CustomerName string             `json:"customer_name"` // Add this field
	Goods        []Goods            `bson:"goods"`
	StartDate    time.Time          `bson:"start_date"`
	EndDate      time.Time          `bson:"end_date"`
	Status       string             `bson:"status"`
	CreatedAt    time.Time          `bson:"created_at"`
}
