package controller

import (
	"encoding/json"
	"net/http"

	"NYANIMEBACKEND/models"
	"NYANIMEBACKEND/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register handler
// Register handler
func Register(w http.ResponseWriter, r *http.Request) {
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
	user.Password = "" // Hapus password dari objek user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user":    user, // Mengembalikan informasi pengguna yang baru terdaftar tanpa password
	})
}

// Login handler
func Login(w http.ResponseWriter, r *http.Request) {
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

	// Return success response with user role
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"user":    user, // Mengembalikan informasi pengguna termasuk role
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Contoh: Hapus token atau session user
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}

// GetAllAnime handler
func GetAllAnime(w http.ResponseWriter, r *http.Request) {
	var animes []models.Anime

	if err := utils.DB.Find(&animes).Error; err != nil {
		http.Error(w, "Failed to fetch anime", http.StatusInternalServerError)
		return
	}

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

	if err := utils.DB.Create(&anime).Error; err != nil {
		http.Error(w, "Failed to create anime", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(anime)
}

// UpdateAnime handler
func EditAnime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var anime models.Anime
	err := json.NewDecoder(r.Body).Decode(&anime)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.DB.Model(&models.Anime{}).Where("id = ?", anime.ID).Updates(anime).Error; err != nil {
		http.Error(w, "Failed to edit anime", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(anime)
}

// DeleteAnime handler
func DeleteAnime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var anime models.Anime
	err := json.NewDecoder(r.Body).Decode(&anime)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := utils.DB.Delete(&anime).Error; err != nil {
		http.Error(w, "Failed to delete anime", http.StatusInternalServerError)
		return
	}

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
