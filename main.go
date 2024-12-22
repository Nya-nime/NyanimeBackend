package main

import (
	"log"
	"net/http"

	"NYANIMEBACKEND/routes"
	"NYANIMEBACKEND/utils"

	"github.com/joho/godotenv" // Import godotenv
)

func main() {
	// Muat variabel lingkungan dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Inisialisasi Database
	utils.InitDB()

	// Setup Routes
	router := routes.SetupRoutes()

	// Jalankan Server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
