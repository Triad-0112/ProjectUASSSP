package services

import (
	"DistributionFlex/models"
	"DistributionFlex/repositories"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ContractService provides business logic for contracts.
type ContractService struct {
	repo         *repositories.ContractRepository
	customerRepo *repositories.CustomerRepository
}

// NewContractService creates a new instance of ContractService.
func NewContractService(repo *repositories.ContractRepository, customerRepo *repositories.CustomerRepository) *ContractService {
	return &ContractService{
		repo:         repo,
		customerRepo: customerRepo,
	}
}

// GetContract retrieves a contract by its ID.
func (s *ContractService) GetContract(ctx context.Context, id string) (*models.Contract, error) {
	return s.repo.FindByID(ctx, id)
}

// CreateContract creates a new contract.
func (s *ContractService) CreateContract(ctx context.Context, contract *models.Contract) (*models.Contract, error) {
	if !contract.CustomerID.IsZero() {
		if _, err := primitive.ObjectIDFromHex(contract.CustomerID.Hex()); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("invalid CustomerID")
	}

	// Ensure the date fields are properly set
	if contract.StartDate.IsZero() || contract.EndDate.IsZero() {
		return nil, errors.New("start date and end date must be set")
	}

	// Set CreatedAt to the current UTC time
	contract.CreatedAt = time.Now().UTC()

	// Save the contract in the repository
	if _, err := s.repo.Create(ctx, contract); err != nil {
		return nil, err
	}

	return contract, nil
}

// UpdateContract updates an existing contract.
func (s *ContractService) UpdateContract(ctx context.Context, id string, update bson.M) (*models.Contract, error) {
	_, err := s.repo.Update(ctx, id, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, id)
}

// GetAllContracts retrieves all contracts from the repository.
func (s *ContractService) GetAllContracts(ctx context.Context) ([]*models.Contract, error) {
	return s.repo.FindAll(ctx)
}

// GetUniqueCustomerIDs extracts unique customer IDs from a list of contracts.
func (s *ContractService) GetUniqueCustomerIDs(contracts []*models.Contract) []primitive.ObjectID {
	return s.repo.GetUniqueCustomerIDs(contracts)
}

// GetContractsWithCustomerNames retrieves all contracts and includes customer names.
func (s *ContractService) GetContractsWithCustomerNames(ctx context.Context) ([]*models.Contract, error) {
	contracts, err := s.GetAllContracts(ctx)
	if err != nil {
		return nil, err
	}

	customerIDs := s.GetUniqueCustomerIDs(contracts)
	customers, err := s.repo.FindCustomersByIDs(ctx, customerIDs) // Adjust according to your implementation
	if err != nil {
		return nil, err
	}

	customerMap := make(map[primitive.ObjectID]*models.Customer)
	for _, customer := range customers {
		customerMap[customer.ID] = customer
	}

	for _, contract := range contracts {
		if customer, ok := customerMap[contract.CustomerID]; ok {
			contract.CustomerName = customer.Name
		}
	}

	return contracts, nil
}
