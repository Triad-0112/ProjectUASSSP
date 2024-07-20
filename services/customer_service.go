package services

import (
	"DistributionFlex/models"
	"DistributionFlex/repositories"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerService struct {
	repo *repositories.CustomerRepository
}

func NewCustomerService(repo *repositories.CustomerRepository) *CustomerService {
	return &CustomerService{
		repo: repo,
	}
}

func (s *CustomerService) GetCustomer(ctx context.Context, id string) (*models.Customer, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CustomerService) CreateCustomer(ctx context.Context, contract *models.Customer) (*models.Customer, error) {
	_, err := s.repo.Create(ctx, contract)
	if err != nil {
		return nil, err
	}
	return contract, nil
}

func (s *CustomerService) UpdateCustomer(ctx context.Context, id string, update bson.M) (*models.Customer, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	updateResult, err := s.repo.Update(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		// No document was matched for the update
		return nil, nil
	}

	updatedCustomer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updatedCustomer, nil
}

func (s *CustomerService) GetAllCustomers(ctx context.Context) ([]*models.Customer, error) {
	return s.repo.FindAll(ctx)
}

func (s *CustomerService) GetCustomersByIDs(ctx context.Context, ids []primitive.ObjectID) ([]*models.Customer, error) {
	// Fetch customers from repository by IDs
	return s.repo.FindByIDs(ctx, ids)
}

func (s *CustomerService) DeleteCustomer(ctx context.Context, id string) error {
	_, err := s.repo.Delete(ctx, id)
	return err
}
