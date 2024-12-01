package main

import (
	"globe-hop/config"
	"globe-hop/models"
	"globe-hop/routes"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("No .env file found, using environment variables")
	}

	// Initialize database
	config.InitDb()
	config.AutoMigrate(&models.User{})

	// Initialize router
	router := routes.InitializeRouter()

	// Start server
	PORT := os.Getenv("PORT_NO")
	log.Printf("Server running on localhost:%v", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
