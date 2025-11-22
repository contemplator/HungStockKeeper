package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Fallback for local development if env var is not set
		// Adjust this connection string to match your local setup
		connStr = "user=postgres password=postgres dbname=hungstockkeeper sslmode=disable"
	}

	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Auto Migrate
	// err = DB.AutoMigrate(&models.User{}, &models.Holding{})
	// if err != nil {
	// 	log.Fatal("Failed to migrate database: ", err)
	// }

	log.Println("Connected to the database")
}
