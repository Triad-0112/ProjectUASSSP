package routes

import (
	"DistributionFlex/controllers"
	"DistributionFlex/middleware"
	"DistributionFlex/repositories"
	"DistributionFlex/services"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewRouter(db *mongo.Database) *mux.Router {
	r := mux.NewRouter()

	userRepository := repositories.NewUserRepository(db)
	contractRepository := repositories.NewContractRepository(db)
	customerRepository := repositories.NewCustomerRepository(db)
	invoiceRepository := repositories.NewInvoiceRepository(db)

	userService := services.NewUserService(userRepository)
	contractService := services.NewContractService(contractRepository, customerRepository)
	customerService := services.NewCustomerService(customerRepository)
	invoiceService := services.NewInvoiceService(invoiceRepository, customerRepository)

	userController := controllers.NewUserController(userService, contractService, customerService, invoiceService)
	customerController := controllers.NewCustomerController(customerService)
	contractController := controllers.NewContractController(contractRepository, customerRepository)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	}).Methods("GET")

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	r.HandleFunc("/login", userController.Login).Methods("POST")
	r.HandleFunc("/logout", userController.Logout).Methods("GET")
	r.Handle("/login", middleware.AuthRequired(http.HandlerFunc(userController.Login))).Methods("GET")         // Protect login page
	r.Handle("/dashboard", middleware.AuthRequired(http.HandlerFunc(userController.Dashboard))).Methods("GET") // Protect dashboard page

	r.Handle("/contract", middleware.AuthRequired(http.HandlerFunc(userController.Contract))).Methods("GET") // Protect contract page
	r.Handle("/contracts", middleware.AuthRequired(http.HandlerFunc(userController.GetContracts))).Methods("GET")

	r.Handle("/customer", middleware.AuthRequired(http.HandlerFunc(userController.Customer))).Methods("GET") // Protect contract page
	r.Handle("/customers", middleware.AuthRequired(http.HandlerFunc(userController.GetCustomer))).Methods("GET")
	r.Handle("/customer/{id}", middleware.AuthRequired(http.HandlerFunc(userController.GetCustomer))).Methods("GET")
	r.Handle("/customer", middleware.AuthRequired(http.HandlerFunc(customerController.CreateCustomer))).Methods("POST")
	r.Handle("/customer/{id}", middleware.AuthRequired(http.HandlerFunc(customerController.UpdateCustomer))).Methods("PUT")
	r.Handle("/customer/{id}", middleware.AuthRequired(http.HandlerFunc(customerController.DeleteCustomer))).Methods("DELETE")

	r.Handle("/invoice", middleware.AuthRequired(http.HandlerFunc(userController.Invoice))).Methods("GET") // Protect contract page
	r.Handle("/invoices", middleware.AuthRequired(http.HandlerFunc(userController.GetInvoice))).Methods("GET")

	r.Handle("/dashboard", middleware.AuthMiddleware(db)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/dashboard.html")
	}))).Methods("GET")

	r.HandleFunc("/contracts", GinHandlerToHTTP(contractController.CreateContract)).Methods("POST")

	return r
}
