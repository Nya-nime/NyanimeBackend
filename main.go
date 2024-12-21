package main

import (
	"log"
	"net/http"

	"NYANIMEBACKEND/routes"
	"NYANIMEBACKEND/utils"
)

func main() {
	// Inisialisasi Database
	utils.InitDB()

	// Setup Routes
	router := routes.SetupRoutes()

	// Jalankan Server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
