package controllers

import (
	"DistributionFlex/services"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
)

type UserController struct {
	service         *services.UserService
	contractService *services.ContractService
	customerService *services.CustomerService
	invoiceService  *services.InvoiceService
}

func NewUserController(service *services.UserService, contractService *services.ContractService, customerService *services.CustomerService, invoiceService *services.InvoiceService) *UserController {
	return &UserController{
		service:         service,
		contractService: contractService,
		customerService: customerService,
		invoiceService:  invoiceService,
	}
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Serve the login page
		cookie, err := r.Cookie("session_token")
		if err == nil && cookie.Value != "" {
			// Redirect to the dashboard if the session token is present
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Unable to load login page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
		return
	}
	if r.Body == nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	user, err := uc.service.Login(context.Background(), username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    user.ID.Hex(), // Use user ID as the session token
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	w.Write([]byte("Login successful"))
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	// Invalidate the session token by setting its MaxAge to -1
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Expire the cookie immediately
	})

	// Redirect to the login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (uc *UserController) Dashboard(w http.ResponseWriter, r *http.Request) {
	// Serve the dashboard page
	tmpl, err := template.ParseFiles("templates/dashboard.html")
	if err != nil {
		http.Error(w, "Unable to load dashboard page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (uc *UserController) Contract(w http.ResponseWriter, r *http.Request) {
	// Serve the contract page
	tmpl, err := template.ParseFiles("templates/contract.html")
	if err != nil {
		http.Error(w, "Unable to load contract page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (uc *UserController) GetContracts(w http.ResponseWriter, r *http.Request) {
	// Fetch all contracts using ContractService
	contracts, err := uc.contractService.GetContractsWithCustomerNames(r.Context())
	if err != nil {
		http.Error(w, "Unable to fetch contracts", http.StatusInternalServerError)
		return
	}
	// Return contracts as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contracts)
}

func (uc *UserController) Customer(w http.ResponseWriter, r *http.Request) {
	// Serve the contract page
	tmpl, err := template.ParseFiles("templates/customer.html")
	if err != nil {
		http.Error(w, "Unable to load contract page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (uc *UserController) GetCustomer(w http.ResponseWriter, r *http.Request) {
	// Log the request for debugging
	log.Println("Received request to fetch customers")

	// Fetch all customers using CustomerService
	customers, err := uc.customerService.GetAllCustomers(r.Context())
	if err != nil {
		log.Printf("Error fetching customers: %v\n", err) // Log the error
		http.Error(w, "Unable to fetch customers", http.StatusInternalServerError)
		return
	}

	log.Println("Fetched customers:", customers) // Log the fetched customers

	// Return customers as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func (uc *UserController) Invoice(w http.ResponseWriter, r *http.Request) {
	// Serve the contract page
	tmpl, err := template.ParseFiles("templates/invoice.html")
	if err != nil {
		http.Error(w, "Unable to load invoice page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (uc *UserController) GetInvoice(w http.ResponseWriter, r *http.Request) {
	// Log the request for debugging
	log.Println("Received request to fetch invoices")

	// Fetch all customers using CustomerService
	invoices, err := uc.invoiceService.GetInvoicessWithCustomerNames(r.Context())
	if err != nil {
		log.Printf("Error fetching invoices: %v\n", err) // Log the error
		http.Error(w, "Unable to fetch invoices", http.StatusInternalServerError)
		return
	}
	// Return customers as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}
