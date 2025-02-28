package controller

import (
	"encoding/json"
	"errors"
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
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if utils.DB == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
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

	log.Printf("User input: %+v\n", user)

	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	var existingUser models.User
	if err := utils.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	user.Role = "user"

	if err := utils.DB.Create(&user).Error; err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

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

// Logout handeler
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

// AddReview handler
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

// CheckUserRating handler
func CheckUserRating(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animeIDStr := vars["anime_id"]
	userIDStr := vars["user_id"]

	animeID, err := strconv.Atoi(animeIDStr)
	if err != nil {
		http.Error(w, "Invalid anime ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var review models.Review
	if err := utils.DB.Where("anime_id = ? AND user_id = ?", animeID, userID).First(&review).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Rating tidak ditemukan, kembalikan status 200 dengan payload kosong
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(nil) // Kembalikan null atau objek kosong
			return
		}
		http.Error(w, "Failed to check rating", http.StatusInternalServerError)
		return
	}

	// Rating ditemukan, kembalikan data rating
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

// EditReview handler
func EditReview(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	vars := mux.Vars(r)              // Ambil variabel dari URL
	reviewIDStr := vars["review_id"] // Ambil ID dari URL

	if reviewIDStr == "" {
		http.Error(w, `{"error": "ID is required"}`, http.StatusBadRequest)
		return
	}

	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid review ID"}`, http.StatusBadRequest)
		return
	}

	var review models.Review
	// Cek apakah review ada
	if err := utils.DB.First(&review, reviewID).Error; err != nil {
		http.Error(w, `{"error": "Review not found"}`, http.StatusNotFound)
		return
	}

	// Decode body request ke dalam updatedReview
	var updatedReview models.Review
	err = json.NewDecoder(r.Body).Decode(&updatedReview)
	if err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validasi data review
	if updatedReview.Content == "" {
		http.Error(w, `{"error": "Content is required"}`, http.StatusBadRequest)
		return
	}
	if updatedReview.Rating < 1 || updatedReview.Rating > 5 {
		http.Error(w, `{"error": "Rating must be between 1 and 5"}`, http.StatusBadRequest)
		return
	}

	// Update field yang ingin diubah
	review.Content = updatedReview.Content
	review.Rating = updatedReview.Rating

	// Simpan perubahan ke database
	if err := utils.DB.Save(&review).Error; err != nil {
		http.Error(w, `{"error": "Failed to update review"}`, http.StatusInternalServerError)
		return
	}

	// Logging
	log.Printf("Editing review with ID: %d", reviewID)
	log.Printf("Updated review data: %+v", review)

	// Set header CORS untuk respons
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(review)
}

// LoadReviews handler
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

// DeleteReview handler
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

// AddFavorite handler
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

// GetFavorites handler
func GetFavorites(w http.ResponseWriter, r *http.Request) {
	// Menangani preflight request untuk CORS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Memastikan metode yang digunakan adalah GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Mengambil userID dari konteks
	userIDValue := r.Context().Value(utils.UserIDKey)
	if userIDValue == nil {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	userID := int(userIDValue.(int))

	// Mengambil daftar favorit dari database dengan join
	var favorites []models.FavoriteWithAnime
	query := `      
		SELECT f.id, f.anime_id, a.title AS anime_title, a.genre, a.description, a.average_rating AS rating, a.release_date      
		FROM favorite f      
		JOIN animes a ON f.anime_id = a.id      
		WHERE f.user_id = ?;      
	`

	if err := utils.DB.Raw(query, userID).Scan(&favorites).Error; err != nil {
		http.Error(w, "Failed to fetch favorites", http.StatusInternalServerError)
		return
	}

	// Log data yang diambil
	log.Printf("Favorites: %+v", favorites)

	// Mengatur header dan mengembalikan respons
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(favorites)
}

// DeleteFavorite handler
func DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	// Menangani preflight request untuk CORS
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Memastikan metode yang digunakan adalah DELETE
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Mengambil anime_id dari URL
	vars := mux.Vars(r)
	favoriteIDStr := vars["id"]
	favoriteID, err := strconv.Atoi(favoriteIDStr)
	if err != nil {
		http.Error(w, "Invalid favorite ID", http.StatusBadRequest)
		return
	}

	// Mengambil userID dari konteks
	userIDValue := r.Context().Value(utils.UserIDKey)
	if userIDValue == nil {
		http.Error(w, "User ID not found", http.StatusUnauthorized)
		return
	}
	userID := int(userIDValue.(int))

	// Menghapus favorite dari database
	if err := utils.DB.Where("id = ? AND user_id = ?", favoriteID, userID).Delete(&models.Favorite{}).Error; err != nil {
		http.Error(w, "Failed to delete favorite", http.StatusInternalServerError)
		return
	}

	// Mengatur header dan mengembalikan respons
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// GetUserProfile data
func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		log.Println("Received OPTIONS request for user profile")
		w.WriteHeader(http.StatusOK)
		return
	}
	userID := r.Context().Value(utils.UserIDKey).(int) // Ambil userID dari konteks

	var user models.User
	if err := utils.DB.Preload("Reviews").First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// GetReviews made by user
func GetUserReviews(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(utils.UserIDKey).(int) // Ambil userID dari konteks

	var reviews []models.ReviewWithAnime // Gunakan struktur ReviewWithAnime

	query := `  
        SELECT r.id, r.anime_id, r.content, r.rating, a.title AS anime_title, a.genre, a.release_date  
        FROM reviews_new r  
        JOIN animes a ON r.anime_id = a.id  
        WHERE r.user_id = ?;  
    `

	if err := utils.DB.Raw(query, userID).Scan(&reviews).Error; err != nil {
		http.Error(w, "No reviews found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reviews)
}

func EditUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(utils.UserIDKey).(int) // Get userID from context

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Find the user by ID
	var user models.User
	if err := utils.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update fields
	user.Username = updatedUser.Username
	user.Bio = updatedUser.Bio

	// Save the updated user
	if err := utils.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
