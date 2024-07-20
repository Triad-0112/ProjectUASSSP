package controllers

import (
	"DistributionFlex/models"
	"DistributionFlex/repositories"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContractController struct {
	contractRepo *repositories.ContractRepository
	customerRepo *repositories.CustomerRepository
}

func NewContractController(contractRepo *repositories.ContractRepository, customerRepo *repositories.CustomerRepository) *ContractController {
	return &ContractController{
		contractRepo: contractRepo,
		customerRepo: customerRepo,
	}
}

func (c *ContractController) GetCustomers(ctx *gin.Context) {
	contracts, err := c.contractRepo.FindAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch customer names
	for i, contract := range contracts {
		customer, err := c.customerRepo.FindByID(ctx, contract.CustomerID.Hex())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if customer != nil {
			contracts[i].CustomerName = customer.Name
			log.Printf("Fetched customer: %v\n", customer) // Add this line for debugging
		} else {
			contracts[i].CustomerName = "Unknown"
		}
	}

	ctx.JSON(http.StatusOK, contracts)
}

func (c *ContractController) CreateContract(ctx *gin.Context) {
	var contract models.Contract
	if err := ctx.BindJSON(&contract); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure CustomerID is properly converted from string to ObjectID
	customerID, err := primitive.ObjectIDFromHex(contract.CustomerID.Hex())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}
	contract.CustomerID = customerID

	// Convert date strings to time.Time
	contract.StartDate, err = time.Parse(time.RFC3339, contract.StartDate.Format(time.RFC3339))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
		return
	}
	contract.EndDate, err = time.Parse(time.RFC3339, contract.EndDate.Format(time.RFC3339))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
		return
	}
	contract.CreatedAt = time.Now().UTC()

	result, err := c.contractRepo.Create(ctx, &contract)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"inserted_id": result.InsertedID})
}
