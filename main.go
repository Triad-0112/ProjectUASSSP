package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"DistributionFlex/routes"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println(err)
		}
	}()

	// Select the database
	db := client.Database(os.Getenv("MONGODB_DB"))

	// Create a new router with MongoDB client and secret key
	router := routes.NewRouter(db)

	// Serve static files from the ./public directory
	fs := http.FileServer(http.Dir("./public"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Serve index.html at the root URL
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	// Retrieve the server port from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	// Start the HTTP server
	log.Println("Server is running on port", port)
	log.Println("Database Used:", os.Getenv("MONGODB_DB"))
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
