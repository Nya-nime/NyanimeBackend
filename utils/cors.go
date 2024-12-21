package utils

import (
	"github.com/rs/cors"
)

// SetupCORS mengatur middleware CORS
func SetupCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"https://nya-nime.github.io/Nyanime/"}, // Ganti dengan domain frontend Anda jika perlu
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
}
