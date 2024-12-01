package config

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb() (*gorm.DB) {
	dsn := os.Getenv("DATABASE_CONNECTION_URL")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Database connection established!")
	return DB
}

func AutoMigrate(models ...interface{}) {
	err := DB.AutoMigrate(models...)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully!")
}