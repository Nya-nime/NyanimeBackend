package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"NYANIMEBACKEND/models"
	"NYANIMEBACKEND/utils"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register handler
func Register(w http.ResponseWriter, r *http.Request) {
	// Menangani permintaan OPTIONS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500") // Ganti dengan domain frontend Anda jika perlu
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi input
	if user.Username == "" || user.Email == "" || user.Password == "" || user.Role == "" {
		http.Error(w, "Username, email, password, and role are required", http.StatusBadRequest)
		return
	}

	// Cek duplikasi email
	var existingUser models.User
	if err := utils.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Save to database
	if err := utils.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Kembalikan respons dengan informasi pengguna tanpa password
	user.Password = ""                                                     // Hapus password dari objek user
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500") // Ganti dengan domain frontend Anda jika perlu
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user":    user, // Mengembalikan informasi pengguna yang baru terdaftar tanpa password
	})
}

// Login handler
func Login(w http.ResponseWriter, r *http.Request) {
	// Menangani permintaan OPTIONS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := utils.DB.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Role) // Pastikan user.ID dan user.Role ada
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return success response with token
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"user":    map[string]interface{}{"id": user.ID, "email": user.Email, "role": user.Role}, // Hanya mengembalikan informasi yang diperlukan
		"token":   token,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Menangani permintaan OPTIONS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500") // Ganti dengan URL frontend Anda
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ambil token dari header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}

	// Token biasanya dalam format "Bearer <token>"
	token := strings.TrimPrefix(authHeader, "Bearer ")
	utils.AddToBlacklist(token)
	// Tambahkan token ke blacklist atau logika lain di sini

	// Kirim respons sukses
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}

// GetAllAnime handler
func GetAllAnime(w http.ResponseWriter, r *http.Request) {
	var animes []models.Anime
	if err := utils.DB.Find(&animes).Error; err != nil {
		http.Error(w, "Failed to retrieve anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(animes)
}

// CreateAnime handler
func CreateAnime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var anime models.Anime
	err := json.NewDecoder(r.Body).Decode(&anime)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi data anime
	if anime.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	log.Printf("Creating anime: %+v", anime)

	if err := utils.DB.Create(&anime).Error; err != nil {
		log.Printf("Error creating anime: %v", err)
		http.Error(w, "Failed to create anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(anime)
}

// UpdateAnime handler
func EditAnime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Ambil variabel dari URL
	id := vars["id"]    // Ambil ID dari URL

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var anime models.Anime
	err := json.NewDecoder(r.Body).Decode(&anime)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi data anime
	if anime.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Update anime di database
	if err := utils.DB.Model(&anime).Where("id = ?", id).Updates(anime).Error; err != nil {
		http.Error(w, "Failed to update anime", http.StatusInternalServerError)
		return
	}

	log.Printf("Editing anime with ID: %s", id)
	log.Printf("Request body: %+v", anime)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(anime)
}

// DeleteAnime handler
func DeleteAnime(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteAnime handler called") // Log untuk debugging

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ambil ID dari URL
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("Attempting to delete anime with ID: %s", id)

	// Cek apakah anime dengan ID tersebut ada
	var anime models.Anime
	if err := utils.DB.First(&anime, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Anime not found", http.StatusNotFound)
			return
		}
		log.Printf("Error finding anime with ID %s: %v", id, err)
		http.Error(w, "Failed to find anime", http.StatusInternalServerError)
		return
	}

	// Hapus anime dari database
	if err := utils.DB.Delete(&anime).Error; err != nil {
		log.Printf("Error deleting anime with ID %s: %v", id, err)
		http.Error(w, "Failed to delete anime", http.StatusInternalServerError)
		return
	}

	if utils.DB.RowsAffected == 0 {
		log.Printf("No anime found with ID: %s", id)
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Anime deleted successfully"})
}

func AddReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value("userID").(int) // Ambil userID dari konteks

	var review models.Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	review.UserID = userID // Set userID untuk review

	if err := utils.DB.Create(&review).Error; err != nil {
		http.Error(w, "Failed to add review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

func EditReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var review models.Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.DB.Save(&review).Error; err != nil {
		http.Error(w, "Failed to update review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(review)
}

func DeleteReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var review models.Review
	err := json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.DB.Delete(&review).Error; err != nil {
		http.Error(w, "Failed to delete review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Review deleted successfully"})
}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int) // Ambil userID dari konteks

	var user models.User
	if err := utils.DB.Preload("Reviews").First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
