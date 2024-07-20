// models/customer.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	Street   string `bson:"street"`
	City     string `bson:"city"`
	Province string `bson:"province"`
	Zip      string `bson:"zip"`
}

type Customer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Phone     string             `bson:"phone"`
	Address   Address            `bson:"address"`
	CreatedAt time.Time          `bson:"created_at"`
}
