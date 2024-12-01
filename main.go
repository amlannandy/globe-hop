package main

import (
	"globe-hop/config"
	"globe-hop/models"
	"globe-hop/routes"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main()  {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("No .env file found, using environment variables")
	}

	// Initialize database
	config.InitDb();
	config.AutoMigrate(&models.User{})


	// Initialize router
	router := routes.InitializeRouter()

	// Start server
	log.Printf("Server running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}