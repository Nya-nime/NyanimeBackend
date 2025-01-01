package utils

import (
	"github.com/rs/cors"
)

// SetupCORS mengatur middleware CORS
func SetupCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5500"}, // Ganti dengan domain frontend Anda jika perlu
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})
}
