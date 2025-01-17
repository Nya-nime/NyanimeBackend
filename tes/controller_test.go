package tes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"NYANIMEBACKEND/controller"
	"NYANIMEBACKEND/models"

	"github.com/gorilla/mux"
)

func setup() {
	InitTestDB()
}

func TestRegister(t *testing.T) {
	setup()
	router := mux.NewRouter()
	router.HandleFunc("/register", controller.Register).Methods("POST", "OPTIONS")

	tests := []struct {
		name       string
		body       models.User
		statusCode int
	}{{
		name: "Valid User",
		body: models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		},
		statusCode: http.StatusCreated,
	}, {
		name: "Missing Email",
		body: models.User{
			Username: "testuser",
			Password: "password123",
		},
		statusCode: http.StatusBadRequest,
	}, {
		name: "Email Already Registered",
		body: models.User{
			Username: "existinguser",
			Email:    "tes@example.com", // Use an email that already exists in the database
			Password: "password123",
		},
		statusCode: http.StatusConflict,
	}}

	// Prepopulate the database with an existing user
	existingUser := models.User{
		Username: "existinguser",
		Email:    "tes@example.com",
		Password: "hashedpassword", // This should be a hashed password
	}
	DB.Create(&existingUser)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	setup()
	router := mux.NewRouter()
	Login := func(w http.ResponseWriter, r *http.Request) {
		// Implement handler logic here
	}
	router.HandleFunc("/login", Login).Methods("POST", "OPTIONS")

	tests := []struct {
		name       string
		body       map[string]string
		statusCode int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

// Add implementation and corrections for other test cases below

func TestLogout(t *testing.T) {
	setup()
	router := mux.NewRouter()
	Logout := func(w http.ResponseWriter, r *http.Request) {
		// Implement handler logic here
	}
	router.HandleFunc("/logout", Logout).Methods("POST", "OPTIONS")

	tests := []struct {
		name       string
		authHeader string
		statusCode int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/logout", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestGetAllAnime(t *testing.T) {
	setup()
	router := mux.NewRouter()
	GetAllAnime := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/anime", GetAllAnime).Methods("GET", "OPTIONS")

	req := httptest.NewRequest("GET", "/anime", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCreateAnime(t *testing.T) {
	setup()
	router := mux.NewRouter()
	CreateAnime := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/anime", CreateAnime).Methods("POST", "OPTIONS")

	tests := []struct {
		name       string
		body       models.Anime
		statusCode int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/anime", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestEditAnime(t *testing.T) {
	setup()
	router := mux.NewRouter()
	EditAnime := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/anime/{id}", EditAnime).Methods("PUT", "OPTIONS")

	tests := []struct {
		name       string
		id         string
		body       models.Anime
		statusCode int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("PUT", "/anime/"+tt.id, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestDeleteAnime(t *testing.T) {
	setup()
	router := mux.NewRouter()
	DeleteAnime := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/anime/{id}", DeleteAnime).Methods("DELETE", "OPTIONS")

	tests := []struct {
		name       string
		id         string
		statusCode int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/anime/"+tt.id, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestAddReview(t *testing.T) {
	setup()
	router := mux.NewRouter()
	AddReview := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/anime/{anime_id}/reviews", AddReview).Methods("POST", "OPTIONS")

	tests := []struct {
		name       string
		animeID    string
		body       models.Review
		statusCode int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/anime/"+tt.animeID+"/reviews", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestCheckUserRating(t *testing.T) {
	setup()
	router := mux.NewRouter()
	CheckUserRating := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/anime/{anime_id}/reviews/{user_id}", CheckUserRating).Methods("GET", "OPTIONS")

	tests := []struct {
		name       string
		animeID    string
		userID     string
		statusCode int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/anime/"+tt.animeID+"/reviews/"+tt.userID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestEditReview(t *testing.T) {
	setup()
	router := mux.NewRouter()
	EditReview := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/reviews/{review_id}", EditReview).Methods("PUT", "OPTIONS")

	tests := []struct {
		name       string
		reviewID   string
		body       models.Review
		statusCode int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("PUT", "/reviews/"+tt.reviewID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestLoadReviews(t *testing.T) {
	setup()
	router := mux.NewRouter()
	LoadReviews := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/anime/{anime_id}/reviews", LoadReviews).Methods("GET", "OPTIONS")

	tests := []struct {
		name       string
		animeID    string
		statusCode int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/anime/"+tt.animeID+"/reviews", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestDeleteReview(t *testing.T) {
	setup()
	router := mux.NewRouter()
	DeleteReview := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/reviews/{review_id}", DeleteReview).Methods("DELETE", "OPTIONS")

	tests := []struct {
		name       string
		reviewID   string
		statusCode int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/reviews/"+tt.reviewID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestAddFavorite(t *testing.T) {
	setup()
	router := mux.NewRouter()
	AddFavorite := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/anime/{anime_id}/favorites", AddFavorite).Methods("POST", "OPTIONS")

	tests := []struct {
		name       string
		animeID    string
		body       models.Favorite
		statusCode int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/anime/"+tt.animeID+"/favorites", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestGetFavorites(t *testing.T) {
	setup()
	router := mux.NewRouter()
	GetFavorites := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/favorites", GetFavorites).Methods("GET", "OPTIONS")

	req := httptest.NewRequest("GET", "/favorites", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestDeleteFavorite(t *testing.T) {
	setup()
	router := mux.NewRouter()
	DeleteFavorite := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/favorites/{id}", DeleteFavorite).Methods("DELETE", "OPTIONS")

	tests := []struct {
		name       string
		favoriteID string
		statusCode int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/favorites/"+tt.favoriteID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestGetUserProfile(t *testing.T) {
	setup()
	router := mux.NewRouter()
	GetUserProfile := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/user/profile", GetUserProfile).Methods("GET", "OPTIONS")

	req := httptest.NewRequest("GET", "/user/profile", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetUserReviews(t *testing.T) {
	setup()
	router := mux.NewRouter()
	GetUserReviews := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/user/reviews", GetUserReviews).Methods("GET", "OPTIONS")

	req := httptest.NewRequest("GET", "/user/reviews", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestEditUserProfile(t *testing.T) {
	setup()
	router := mux.NewRouter()
	EditUserProfile := func(w http.ResponseWriter, r *http.Request) {
		// Handler logic here
	}
	router.HandleFunc("/user/profile", EditUserProfile).Methods("PUT", "OPTIONS")

	tests := []struct {
		name       string
		body       models.User
		statusCode int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("PUT", "/user/profile", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
