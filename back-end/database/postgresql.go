package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/louismomo66/logger/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v\n", err)
		return nil, err
	}
	databaseURL := os.Getenv("DATABASE_URL")
	// dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
	// 	"db",       // Use the service name defined in docker-compose.yml
	// 	5432,       // Standard PostgreSQL port
	// 	"devicel",  // Username from docker-compose.yml
	// 	"postres1", // Password from docker-compose.yml
	// 	"logger",   // Database name from docker-compose.yml
	// )
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  databaseURL, // Add createDatabase=true parameter
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	// conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return conn, err
	}

	// Auto-migrate models
	err = conn.AutoMigrate(&models.Device{}, &models.Readings{},&models.User{})
	if err != nil {
		log.Printf("Failed to auto-migrate: %v", err)
		return nil, err
	}

	return conn, nil
}
