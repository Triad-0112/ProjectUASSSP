package services

import (
	"DistributionFlex/models"
	"DistributionFlex/repositories"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InvoiceService struct {
	repo         *repositories.InvoiceRepository
	customerRepo *repositories.CustomerRepository
}

func NewInvoiceService(repo *repositories.InvoiceRepository, customerRepo *repositories.CustomerRepository) *InvoiceService {
	return &InvoiceService{
		repo:         repo,
		customerRepo: customerRepo,
	}
}

func (s *InvoiceService) GetInvoice(ctx context.Context, id string) (*models.Invoice, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *InvoiceService) CreateInvoice(ctx context.Context, contract *models.Invoice) (*models.Invoice, error) {
	_, err := s.repo.Create(ctx, contract)
	if err != nil {
		return nil, err
	}
	return contract, nil
}

func (s *InvoiceService) GetAllInvoices(ctx context.Context) ([]*models.Invoice, error) {
	return s.repo.FindAll(ctx)
}

func (s *InvoiceService) GetUniqueCustomerIDs(invoices []*models.Invoice) []primitive.ObjectID {
	return s.repo.GetUniqueCustomerIDs(invoices)
}

func (s *InvoiceService) GetInvoicessWithCustomerNames(ctx context.Context) ([]*models.Invoice, error) {
	invoices, err := s.GetAllInvoices(ctx)
	if err != nil {
		return nil, err
	}

	customerIDs := s.GetUniqueCustomerIDs(invoices)
	customers, err := s.repo.FindCustomersByIDs(ctx, customerIDs) // Adjust according to your implementation
	if err != nil {
		return nil, err
	}

	customerMap := make(map[primitive.ObjectID]*models.Customer)
	for _, customer := range customers {
		customerMap[customer.ID] = customer
	}

	for _, invoice := range invoices {
		if customer, ok := customerMap[invoice.CustomerID]; ok {
			invoice.CustomerName = customer.Name
		}
	}

	return invoices, nil
}
