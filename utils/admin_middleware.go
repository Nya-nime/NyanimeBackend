package utils

import (
	"log"
	"net/http"
)

// Middleware untuk otorisasi admin
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(UserRoleKey).(string)

		// Tambahkan log untuk debugging
		log.Printf("User role: %s", role)

		// Periksa apakah role ada dan apakah itu admin
		if !ok || role != "admin" {
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
