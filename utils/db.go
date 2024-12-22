package utils

import (
	"NYANIMEBACKEND/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // Import godotenv
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Memuat file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Ambil variabel lingkungan (environment variables)
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Debugging: Tampilkan nilai variabel lingkungan untuk memastikan mereka benar
	fmt.Println("DB_USER:", dbUser)
	fmt.Println("DB_PASS:", dbPass)
	fmt.Println("DB_HOST:", dbHost)
	fmt.Println("DB_PORT:", dbPort)
	fmt.Println("DB_NAME:", dbName)

	// Periksa jika ada variabel lingkungan yang kosong
	if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatalf("One or more environment variables are missing. Please check DB_USER, DB_PASS, DB_HOST, DB_PORT, DB_NAME")
	}

	// Format DSN (Data Source Name) untuk MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	// Koneksi ke database
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Logging jika berhasil terkoneksi
	log.Println("Database connected successfully!")

	// Opsional: Auto-migrate model
	autoMigrateModels()
}

func autoMigrateModels() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Anime{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}
	log.Println("Database models migrated successfully!")
}
