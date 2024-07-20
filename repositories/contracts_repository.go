package repositories

import (
	"DistributionFlex/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ContractRepository handles database operations for contracts.
type ContractRepository struct {
	collection *mongo.Collection
}

// NewContractRepository creates a new instance of ContractRepository.
func NewContractRepository(db *mongo.Database) *ContractRepository {
	return &ContractRepository{
		collection: db.Collection("contracts"),
	}
}

// FindByID retrieves a contract by its ID.
func (r *ContractRepository) FindByID(ctx context.Context, id string) (*models.Contract, error) {
	var contract models.Contract
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&contract)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &contract, nil
}

// Create inserts a new contract into the collection.
func (r *ContractRepository) Create(ctx context.Context, contract *models.Contract) (*mongo.InsertOneResult, error) {
	return r.collection.InsertOne(ctx, contract)
}

// Update updates an existing contract.
func (r *ContractRepository) Update(ctx context.Context, id string, update bson.M) (*mongo.UpdateResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
}

func (r *ContractRepository) FindAll(ctx context.Context) ([]*models.Contract, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var contracts []*models.Contract
	for cursor.Next(ctx) {
		var contract models.Contract
		if err := cursor.Decode(&contract); err != nil {
			return nil, err
		}
		contracts = append(contracts, &contract)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return contracts, nil
}

func (r *ContractRepository) GetUniqueCustomerIDs(contracts []*models.Contract) []primitive.ObjectID {
	customerIDs := make(map[primitive.ObjectID]struct{})
	for _, contract := range contracts {
		customerIDs[contract.CustomerID] = struct{}{}
	}

	var ids []primitive.ObjectID
	for id := range customerIDs {
		ids = append(ids, id)
	}
	return ids
}

func (r *ContractRepository) FindCustomersByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*models.Customer, error) {
	cursor, err := r.collection.Database().Collection("customers").Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
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
