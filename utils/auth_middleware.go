package utils

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Key untuk context
type ContextKey string

const UserIDKey ContextKey = "userID"
const UserRoleKey ContextKey = "userRole"

// Middleware untuk autentikasi
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Verifikasi token
		token, claims, err := VerifyToken(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Simpan informasi user ke context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// VerifyToken memverifikasi JWT dan mengembalikan klaim
func VerifyToken(tokenString string) (*jwt.Token, *CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("G3n3r@t3dS3cr3tK3y!2024"), nil
	})
	return token, claims, err
}

// CustomClaims untuk JWT
type CustomClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Blacklist untuk menyimpan token yang diblacklist
var tokenBlacklist = struct {
	sync.RWMutex
	m map[string]time.Time
}{m: make(map[string]time.Time)}

// AddToBlacklist menambahkan token ke blacklist
func AddToBlacklist(token string) {
	tokenBlacklist.Lock()
	defer tokenBlacklist.Unlock()
	tokenBlacklist.m[token] = time.Now() // Menyimpan waktu saat token diblacklist
}

// IsBlacklisted memeriksa apakah token ada di blacklist
func IsBlacklisted(token string) bool {
	tokenBlacklist.RLock()
	defer tokenBlacklist.RUnlock()
	_, exists := tokenBlacklist.m[token]
	return exists
}
