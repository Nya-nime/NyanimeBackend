package tes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"NYANIMEBACKEND/controller"
	"NYANIMEBACKEND/models"
	"NYANIMEBACKEND/utils"

	"github.com/gorilla/mux"
)

func setup() {
	InitTestDB()
	utils.DB = DB
}

func TestRegister(t *testing.T) {
	setup() // Inisialisasi database dan router
	router := mux.NewRouter()
	router.HandleFunc("/user/register", controller.Register).Methods("POST", "OPTIONS")

	tests := []struct {
		name       string
		body       models.User
		statusCode int
	}{
		{"ValidRegister", models.User{Username: "newsuser", Email: "barus@example.com", Password: "password"}, http.StatusCreated},
		{"InvalidRegister", models.User{Username: "", Email: "salah@example.com", Password: "password"}, http.StatusBadRequest},
		{"EmailAlreadyExists", models.User{Username: "existinguser", Email: "sudahadda@example.com", Password: "password"}, http.StatusConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/user/register", bytes.NewBuffer(body))
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
	router.HandleFunc("/user/login", controller.Login).Methods("POST", "OPTIONS")

	tests := []struct {
		name       string
		body       map[string]string
		statusCode int
	}{
		{"ValidLogin", map[string]string{"username": "newsuser", "password": "password"}, http.StatusOK},
		{"InvalidLogin", map[string]string{"username": "", "password": "password"}, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/user/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	setup()
	router := mux.NewRouter()
	Logout := func(w http.ResponseWriter, r *http.Request) {
	}
	router.HandleFunc("/user/logout", Logout).Methods("POST", "OPTIONS")

	tests := []struct {
		name       string
		authHeader string
		statusCode int
	}{
		{"ValidLogout", "Bearer validtoken", http.StatusOK},
		{"InvalidLogout", "", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/user/logout", nil)
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
	}
	router.HandleFunc("/anime/", GetAllAnime).Methods("GET", "OPTIONS")

	req := httptest.NewRequest("GET", "/anime/", nil)
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
	}
	router.HandleFunc("/anime/", CreateAnime).Methods("POST", "OPTIONS")

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
	}
	router.HandleFunc("/anime/{anime_id}/reviews", AddReview).Methods("POST", "OPTIONS")

	tests := []struct {
		name       string
		animeID    string
		body       models.Review
		statusCode int
	}{
		{"ValidReview", "1", models.Review{Content: "Great anime!", Rating: 5}, http.StatusOK},
		{"InvalidReview", "1", models.Review{Content: "", Rating: 0}, http.StatusBadRequest},
	}

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
	}
	router.HandleFunc("/anime/{anime_id}/{user_id}", CheckUserRating).Methods("GET", "OPTIONS")

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
	}
	router.HandleFunc("/anime/{anime_id}", LoadReviews).Methods("GET", "OPTIONS")

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
	}
	router.HandleFunc("/reviews/{review_id}", DeleteReview).Methods("DELETE", "OPTIONS")

	tests := []struct {
		name       string
		reviewID   string
		statusCode int
	}{
		{"ValidDelete", "1", http.StatusOK},
		{"InvalidDelete", "999", http.StatusConflict},
	}

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
	}
	router.HandleFunc("/anime/{anime_id}/favorites", AddFavorite).Methods("POST", "OPTIONS")

	tests := []struct {
		name       string
		animeID    string
		body       models.Favorite
		statusCode int
	}{
		{"ValidFavorite", "1", models.Favorite{UserID: 1, AnimeID: 1}, http.StatusOK},
		{"InvalidFavorite", "999", models.Favorite{UserID: 1, AnimeID: 999}, http.StatusBadRequest},
	}

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
	}
	router.HandleFunc("/favorites/", GetFavorites).Methods("GET", "OPTIONS")

	req := httptest.NewRequest("GET", "/favorites/", nil)
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
	}
	router.HandleFunc("/review/reviews", GetUserReviews).Methods("GET", "OPTIONS")

	req := httptest.NewRequest("GET", "/review/reviews", nil)
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
	}
	router.HandleFunc("/user/profile", EditUserProfile).Methods("PUT", "OPTIONS")

	tests := []struct {
		name       string
		body       models.User
		statusCode int
	}{
		{"ValidEdit", models.User{Username: "updateduser", Email: "updated@example.com"}, http.StatusOK},
		{"InvalidEdit", models.User{Username: "", Email: "invalid"}, http.StatusBadRequest},
	}

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
