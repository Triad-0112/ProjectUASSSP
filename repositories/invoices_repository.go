package repositories

import (
	"DistributionFlex/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceRepository struct {
	collection *mongo.Collection
}

func NewInvoiceRepository(db *mongo.Database) *InvoiceRepository {
	return &InvoiceRepository{
		collection: db.Collection("invoices"),
	}
}

func (r *InvoiceRepository) Create(ctx context.Context, Invoice *models.Invoice) (*mongo.InsertOneResult, error) {
	return r.collection.InsertOne(ctx, Invoice)
}

func (r *InvoiceRepository) Update(ctx context.Context, id string, update bson.M) (*mongo.UpdateResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	return r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
}

func (r *InvoiceRepository) FindByID(ctx context.Context, id string) (*models.Invoice, error) {
	var Invoice models.Invoice
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&Invoice)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &Invoice, nil
}

func (r *InvoiceRepository) FindAll(ctx context.Context) ([]*models.Invoice, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var Invoices []*models.Invoice
	for cursor.Next(ctx) {
		var Invoice models.Invoice
		if err := cursor.Decode(&Invoice); err != nil {
			return nil, err
		}
		Invoices = append(Invoices, &Invoice)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return Invoices, nil
}

func (r *InvoiceRepository) GetUniqueCustomerIDs(invoices []*models.Invoice) []primitive.ObjectID {
	customerIDs := make(map[primitive.ObjectID]struct{})
	for _, invoice := range invoices {
		customerIDs[invoice.CustomerID] = struct{}{}
	}

	var ids []primitive.ObjectID
	for id := range customerIDs {
		ids = append(ids, id)
	}
	return ids
}

func (r *InvoiceRepository) FindCustomersByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*models.Customer, error) {
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
