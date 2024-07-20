package repositories

import (
	"DistributionFlex/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepository struct {
	collection *mongo.Collection
}

func NewCustomerRepository(db *mongo.Database) *CustomerRepository {
	return &CustomerRepository{
		collection: db.Collection("customers"),
	}
}

func (r *CustomerRepository) Create(ctx context.Context, customer *models.Customer) (*mongo.InsertOneResult, error) {
	return r.collection.InsertOne(ctx, customer)
}

func (r *CustomerRepository) Update(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return r.collection.UpdateOne(ctx, filter, update)
}

func (r *CustomerRepository) FindByID(ctx context.Context, id string) (*models.Customer, error) {
	var customer models.Customer
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepository) FindAll(ctx context.Context) ([]*models.Customer, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var customers []*models.Customer
	for cursor.Next(ctx) {
		var customer models.Customer
		if err := cursor.Decode(&customer); err != nil {
			return nil, err
		}
		customers = append(customers, &customer)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *CustomerRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*models.Customer, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var customers []*models.Customer
	for cursor.Next(ctx) {
		var customer models.Customer
		if err := cursor.Decode(&customer); err != nil {
			return nil, err
		}
		customers = append(customers, &customer)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *CustomerRepository) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
}
