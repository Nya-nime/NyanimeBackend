package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validasi input tanpa role
	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
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

	// Set role default jika diperlukan (misalnya, "user")
	user.Role = "user" // Atur role default jika Anda ingin semua pengguna terdaftar sebagai user

	// Save to database
	if err := utils.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Kembalikan respons dengan informasi pengguna tanpa password
	user.Password = ""
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user":    user,
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
	// Handle OPTIONS request for CORS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500") // Update with your frontend URL
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is required", http.StatusUnauthorized)
		return
	}

	// Token is usually in the format "Bearer <token>"
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Add the token to the blacklist
	utils.AddToBlacklist(token)

	// Send success response
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Content-Type", "text/plain") // Set content type
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}

// GetAllAnime handler
func GetAllAnime(w http.ResponseWriter, r *http.Request) {
	// Cek metode OPTIONS untuk CORS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	var animes []models.Anime
	// Query untuk mengambil anime dengan rating akumulatif dari tabel animes
	if err := utils.DB.Table("animes"). // Ganti "anime" dengan "animes"
						Select("animes.*, COALESCE(AVG(reviews_new.rating), 0) as average_rating"). // Ganti "reviews" dengan "reviews_new"
						Joins("LEFT JOIN reviews_new ON animes.id = reviews_new.anime_id").         // Ganti "reviews" dengan "reviews_new"
						Group("animes.id").                                                         // Ganti "anime" dengan "animes"
						Scan(&animes).Error; err != nil {
		log.Println("Error retrieving anime:", err) // Log error untuk debugging
		http.Error(w, "Failed to retrieve anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(animes)
}

// CreateAnime handler
func CreateAnime(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}
	log.Println("CreateAnime called") // Log untuk memastikan fungsi dipanggil
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var anime models.Anime
	err := json.NewDecoder(r.Body).Decode(&anime)
	if err != nil {
		log.Printf("Error decoding request body: %v", err) // Log kesalahan decoding
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Log data yang diterima
	log.Printf("Received anime data: %+v", anime)

	// Validasi data anime
	if anime.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	// Log sebelum menyimpan ke database
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
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

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
	// Menangani permintaan OPTIONS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500") // Ganti dengan domain frontend Anda
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	var anime models.Anime
	if err := utils.DB.First(&anime, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Anime not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to find anime", http.StatusInternalServerError)
		return
	}

	if err := utils.DB.Delete(&anime).Error; err != nil {
		http.Error(w, "Failed to delete anime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Anime deleted successfully"})
}

func AddReview(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	animeIDStr := vars["anime_id"]
	animeID, err := strconv.Atoi(animeIDStr)
	if err != nil {
		http.Error(w, "Invalid anime ID", http.StatusBadRequest)
		return
	}

	var review models.Review
	err = json.NewDecoder(r.Body).Decode(&review)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Cek apakah anime_id valid
	var anime models.Anime
	if err := utils.DB.First(&anime, animeID).Error; err != nil {
		http.Error(w, "Anime not found", http.StatusBadRequest)
		return
	}

	// Set userID dan animeID untuk review
	userIDValue := r.Context().Value(utils.UserIDKey)
	if userIDValue != nil {
		review.UserID = int(userIDValue.(int))
	} else {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	review.AnimeID = uint(animeID)

	if err := utils.DB.Create(&review).Error; err != nil {
		http.Error(w, "Failed to create review", http.StatusInternalServerError)
		return
	}

	var averageRating float64
	if err := utils.DB.Table("reviews_new").
		Select("COALESCE(AVG(rating), 0)").
		Where("anime_id = ?", animeID).
		Scan(&averageRating).Error; err != nil {
		log.Println("Error calculating average rating:", err)
	} else {
		// Update average rating di tabel animes
		if err := utils.DB.Model(&anime).Update("average_rating", averageRating).Error; err != nil {
			log.Println("Error updating average rating:", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

func EditReview(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	reviewIDStr := vars["review_id"]
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	var review models.Review
	// Cek apakah review ada
	if err := utils.DB.First(&review, reviewID).Error; err != nil {
		http.Error(w, "Review not found", http.StatusNotFound)
		return
	}

	// Decode body request ke dalam review
	var updatedReview models.Review
	err = json.NewDecoder(r.Body).Decode(&updatedReview)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update field yang ingin diubah
	review.Content = updatedReview.Content
	review.Rating = updatedReview.Rating

	// Simpan perubahan ke database
	if err := utils.DB.Save(&review).Error; err != nil {
		http.Error(w, "Failed to update review", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func LoadReviews(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	animeIDStr := vars["anime_id"]
	animeID, err := strconv.Atoi(animeIDStr)
	if err != nil {
		http.Error(w, "Invalid anime ID", http.StatusBadRequest)
		return
	}

	var reviews []models.Review
	if err := utils.DB.Where("anime_id = ?", animeID).Find(&reviews).Error; err != nil {
		http.Error(w, "Failed to load reviews", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func DeleteReview(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	reviewIDStr := vars["review_id"]
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	// Hapus review
	if err := utils.DB.Delete(&models.Review{}, reviewID).Error; err != nil {
		http.Error(w, "Failed to delete review", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func AddFavorite(w http.ResponseWriter, r *http.Request) {
	// Menangani preflight request untuk CORS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Memastikan metode yang digunakan adalah POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Mengambil anime_id dari URL
	vars := mux.Vars(r)
	animeIDStr := vars["anime_id"]
	animeID, err := strconv.Atoi(animeIDStr)
	if err != nil {
		http.Error(w, "Invalid anime ID", http.StatusBadRequest)
		return
	}

	// Membaca body request untuk favorite
	var favorite models.Favorite
	if err := json.NewDecoder(r.Body).Decode(&favorite); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Cek apakah anime_id valid
	var anime models.Anime
	if err := utils.DB.First(&anime, animeID).Error; err != nil {
		http.Error(w, "Anime not found", http.StatusNotFound)
		return
	}

	// Set userID dan animeID untuk favorite
	userIDValue := r.Context().Value(utils.UserIDKey)
	if userIDValue != nil {
		favorite.UserID = int(userIDValue.(int)) // Mengambil userID dari konteks
	} else {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}

	favorite.AnimeID = uint64(animeID) // Pastikan AnimeID adalah uint64
	favorite.AnimeTitle = anime.Title

	// Simpan ke database
	if err := utils.DB.Create(&favorite).Error; err != nil {
		http.Error(w, "Failed to add favorite", http.StatusInternalServerError)
		return
	}

	// Mengatur header dan mengembalikan respons
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(favorite)
}

func GetFavorites(w http.ResponseWriter, r *http.Request) {
	userIDValue := r.Context().Value("userID")
	if userIDValue == nil {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var favorites []models.Favorite
	if err := utils.DB.Where("user_id = ?", userID).Find(&favorites).Error; err != nil {
		http.Error(w, "Failed to load favorites", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favorites)
}

func RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	favoriteIDStr := vars["id"]
	favoriteID, err := strconv.Atoi(favoriteIDStr)
	if err != nil {
		http.Error(w, "Invalid favorite ID", http.StatusBadRequest)
		return
	}

	// Hapus favorit dari database
	if err := utils.DB.Delete(&models.Favorite{}, favoriteID).Error; err != nil {
		http.Error(w, "Failed to delete favorite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey).(int) // Ambil userID dari konteks

	var user models.User
	if err := utils.DB.Preload("Reviews").First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
