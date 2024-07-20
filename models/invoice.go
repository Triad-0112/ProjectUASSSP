package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ContractID   primitive.ObjectID `bson:"contract_id"`
	CustomerName string             `json:"customer_name"`
	CustomerID   primitive.ObjectID `bson:"customer_id"`
	Status       string             `bson:"status"`
	Amount       int                `bson:"amount"`
	Goods        []Goods            `bson:"goods"`
	InvoiceDate  time.Time          `bson:"invoice_date"`
	DueDate      time.Time          `bson:"due_date"`
	CreatedAt    time.Time          `bson:"created_at"`
}
