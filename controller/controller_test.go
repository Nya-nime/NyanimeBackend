package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	utils "NYANIMEBACKEND/config"
	"NYANIMEBACKEND/models"

	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	// Inisialisasi database untuk pengujian
	return utils.InitTestDB() // Pastikan Anda memiliki fungsi ini untuk menginisialisasi DB
}

func TestRegister(t *testing.T) {
	// Inisialisasi database untuk pengujian
	_ = utils.InitTestDB() // Pastikan Anda memanggil InitTestDB

	// Skenario 1: Pendaftaran pengguna yang berhasil
	t.Run("Successful Registration", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
			Role:     "user",
		}
		body, _ := json.Marshal(user)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Register)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var response map[string]interface{}
		json.NewDecoder(rr.Body).Decode(&response)
		if response["message"] != "User registered successfully" {
			t.Errorf("handler returned unexpected body: got %v want %v", response["message"], "User registered successfully")
		}
	})

	// Tambahkan skenario lain di sini...
}
func TestLogin(t *testing.T) {
	_ = setupTestDB() // Inisialisasi database, tidak perlu menyimpan ke variabel

	// Register a user first
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user",
	}
	utils.DB.Create(&user)

	loginRequest := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(loginRequest)

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetAllAnime(t *testing.T) {
	_ = setupTestDB() // Inisialisasi database, tidak perlu menyimpan ke variabel

	req, err := http.NewRequest("GET", "/anime", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetAllAnime)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestCreateAnime(t *testing.T) {
	_ = setupTestDB() // Inisialisasi database, tidak perlu menyimpan ke variabel

	anime := models.Anime{
		Title: "Test Anime",
	}

	body, _ := json.Marshal(anime)

	req, err := http.NewRequest("POST", "/anime", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateAnime)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestEditAnime(t *testing.T) {
	_ = setupTestDB() // Inisialisasi database, tidak perlu menyimpan ke variabel

	// Create an anime first
	anime := models.Anime{
		Title: "Test Anime",
	}
	utils.DB.Create(&anime)

	updatedAnime := models.Anime{
		Title: "Updated Anime",
	}
	body, _ := json.Marshal(updatedAnime)

	req, err := http.NewRequest("PUT", "/anime/"+strconv.Itoa(int(anime.ID)), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(EditAnime)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeleteAnime(t *testing.T) {
	_ = setupTestDB() // Inisialisasi database, tidak perlu menyimpan ke variabel

	// Create an anime first
	anime := models.Anime{
		Title: "Test Anime",
	}
	utils.DB.Create(&anime)

	req, err := http.NewRequest("DELETE", "/anime/"+strconv.Itoa(int(anime.ID)), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteAnime)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAddReview(t *testing.T) {
	_ = setupTestDB() // Inisialisasi database, tidak perlu menyimpan ke variabel

	// Create an anime first
	anime := models.Anime{
		Title: "Test Anime",
	}
	utils.DB.Create(&anime)

	review := models.Review{
		Content: "Great anime!",
		Rating:  5,
	}

	body, _ := json.Marshal(review)

	req, err := http.NewRequest("POST", "/anime/"+strconv.Itoa(int(anime.ID))+"/reviews", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddReview)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestLoadReviews(t *testing.T) {
	_ = setupTestDB() // Inisialisasi database, tidak perlu menyimpan ke variabel

	// Create an anime and a review first
	anime := models.Anime{
		Title: "Test Anime",
	}
	utils.DB.Create(&anime)

	review := models.Review{
		Content: "Great anime!",
		Rating:  5,
		AnimeID: anime.ID,
	}
	utils.DB.Create(&review)

	req, err := http.NewRequest("GET", "/anime/"+strconv.Itoa(int(anime.ID))+"/reviews", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoadReviews)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeleteReview(t *testing.T) {
	_ = setupTestDB() // Inisialisasi database, tidak perlu menyimpan ke variabel

	// Create an anime and a review first
	anime := models.Anime{
		Title: "Test Anime",
	}
	utils.DB.Create(&anime)

	review := models.Review{
		Content: "Great anime!",
		Rating:  5,
		AnimeID: anime.ID,
	}
	utils.DB.Create(&review)

	req, err := http.NewRequest("DELETE", "/reviews/"+strconv.Itoa(int(review.ID)), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteReview)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestGetUserProfile(t *testing.T) {
	_ = setupTestDB() // Inisialisasi database, tidak perlu menyimpan ke variabel

	// Create a user first
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     "user",
	}
	utils.DB.Create(&user)

	req, err := http.NewRequest("GET", "/user/profile", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUserProfile)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
