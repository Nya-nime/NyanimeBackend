package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken menghasilkan JWT untuk pengguna
func GenerateToken(userID int, role string) (string, error) {
	claims := &CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token berlaku selama 24 jam
			Issuer:    "Nyanime",                                          // Ganti dengan nama aplikasi Anda
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("G3n3r@t3dS3cr3tK3y!2024")) // Ganti dengan kunci rahasia Anda
}
