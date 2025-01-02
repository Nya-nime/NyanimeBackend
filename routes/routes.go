package routes

import (
	"NYANIMEBACKEND/controller"
	"NYANIMEBACKEND/utils"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	corsHandler := utils.SetupCORS()

	// Apply CORS middleware to the router
	router.Use(corsHandler.Handler)

	router.HandleFunc("/user/register", controller.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/login", controller.Login).Methods("POST", "OPTIONS")

	// Redirect root to the frontend
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://127.0.0.1:5500", http.StatusFound)
	}).Methods("GET")

	// User Routes
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/register", controller.Register).Methods("OPTIONS", "POST")
	userRouter.HandleFunc("/login", controller.Login).Methods("OPTIONS", "POST")
	userRouter.Handle("/profile", utils.AuthMiddleware(http.HandlerFunc(controller.GetUserProfile))).Methods("GET")
	userRouter.Handle("/logout", utils.AuthMiddleware(http.HandlerFunc(controller.Logout))).Methods("OPTIONS", "POST")

	// Anime Routes (Admin Privileges)
	animeRouter := router.PathPrefix("/anime").Subrouter()
	animeRouter.HandleFunc("/", controller.GetAllAnime).Methods("GET", "OPTIONS")                                                                   // Semua user bisa lihat
	animeRouter.Handle("/", utils.AuthMiddleware(utils.AdminMiddleware(http.HandlerFunc(controller.CreateAnime)))).Methods("OPTIONS", "POST")       // Menambahkan anime
	animeRouter.Handle("/{id}", utils.AuthMiddleware(utils.AdminMiddleware(http.HandlerFunc(controller.EditAnime)))).Methods("OPTIONS", "PUT")      // Mengedit anime
	animeRouter.Handle("/{id}", utils.AuthMiddleware(utils.AdminMiddleware(http.HandlerFunc(controller.DeleteAnime)))).Methods("OPTIONS", "DELETE") // Menghapus anime

	// Rute untuk menyajikan halaman admin
	router.HandleFunc("/admin.html", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Jika tidak ada token, arahkan ke halaman login
			http.Redirect(w, r, "/login.html", http.StatusFound)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, _, err := utils.VerifyToken(tokenString)
		if err != nil || !token.Valid {
			// Jika token tidak valid, arahkan ke halaman login
			http.Redirect(w, r, "/llogin.html", http.StatusFound)
			return
		}

		// Jika token valid, sajikan halaman admin
		http.ServeFile(w, r, "path/to/admin_home.html") // Ganti dengan path yang sesuai
	}).Methods("GET")

	// Rute untuk menyajikan halaman user_home.html
	router.HandleFunc("/user_home.html", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Jika tidak ada token, arahkan ke halaman login
			http.Redirect(w, r, "/login.html", http.StatusFound)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, _, err := utils.VerifyToken(tokenString)
		if err != nil || !token.Valid {
			// Jika token tidak valid, arahkan ke halaman login
			http.Redirect(w, r, "/login.html", http.StatusFound)
			return
		}

		// Jika token valid, sajikan halaman user_home.html
		http.ServeFile(w, r, "path/to/user_home.html") // Ganti dengan path yang sesuai
	}).Methods("GET")

	// Rute untuk menyajikan halaman profile.html
	router.HandleFunc("/profile.html", func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// Jika tidak ada token, arahkan ke halaman login
			http.Redirect(w, r, "/login.html", http.StatusFound)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, _, err := utils.VerifyToken(tokenString)
		if err != nil || !token.Valid {
			// Jika token tidak valid, arahkan ke halaman login
			http.Redirect(w, r, "/login.html", http.StatusFound)
			return
		}

		// Jika token valid, sajikan halaman profile.html
		http.ServeFile(w, r, "path/to/profile.html") // Ganti dengan path yang sesuai
	}).Methods("GET")

	// Review Routes
	reviewRouter := router.PathPrefix("/review").Subrouter()
	reviewRouter.Handle("/anime/{anime_id}", utils.AuthMiddleware(http.HandlerFunc(controller.AddReview))).Methods("POST")
	reviewRouter.Handle("/{id}", utils.AuthMiddleware(http.HandlerFunc(controller.EditReview))).Methods("PUT")
	reviewRouter.Handle("/{id}", utils.AuthMiddleware(http.HandlerFunc(controller.DeleteReview))).Methods("DELETE")

	return router
}
