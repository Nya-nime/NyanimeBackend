package routes

import (
	"NYANIMEBACKEND/controller"
	"NYANIMEBACKEND/utils"
	"net/http"

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
	}).Methods("GET", "OPTIONS")

	// User Routes
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/register", controller.Register).Methods("OPTIONS", "POST")
	userRouter.HandleFunc("/login", controller.Login).Methods("OPTIONS", "POST")
	userRouter.Handle("/profile", utils.AuthMiddleware(http.HandlerFunc(controller.GetUserProfile))).Methods("GET", "OPTIONS")
	userRouter.Handle("/logout", utils.AuthMiddleware(http.HandlerFunc(controller.Logout))).Methods("OPTIONS", "POST")

	// Anime Routes (Admin Privileges)
	animeRouter := router.PathPrefix("/anime").Subrouter()
	animeRouter.HandleFunc("/", controller.GetAllAnime).Methods("GET", "OPTIONS")                                                                   // Semua user bisa lihat
	animeRouter.Handle("/", utils.AuthMiddleware(utils.AdminMiddleware(http.HandlerFunc(controller.CreateAnime)))).Methods("OPTIONS", "POST")       // Menambahkan anime
	animeRouter.Handle("/{id}", utils.AuthMiddleware(utils.AdminMiddleware(http.HandlerFunc(controller.EditAnime)))).Methods("OPTIONS", "PUT")      // Mengedit anime
	animeRouter.Handle("/{id}", utils.AuthMiddleware(utils.AdminMiddleware(http.HandlerFunc(controller.DeleteAnime)))).Methods("OPTIONS", "DELETE") // Menghapus anime

	// Review Routes
	reviewRouter := router.PathPrefix("/review").Subrouter()
	reviewRouter.Handle("/anime/{anime_id}", utils.AuthMiddleware(http.HandlerFunc(controller.LoadReviews))).Methods("GET")
	reviewRouter.Handle("/anime/{anime_id}/{user_id}", utils.AuthMiddleware(http.HandlerFunc(controller.CheckUserRating))).Methods("OPTIONS", "GET")
	reviewRouter.Handle("/reviews", utils.AuthMiddleware(http.HandlerFunc(controller.GetUserReviews))).Methods("GET", "OPTIONS")
	reviewRouter.Handle("/anime/{anime_id}", utils.AuthMiddleware(http.HandlerFunc(controller.AddReview))).Methods("OPTIONS", "POST")
	reviewRouter.Handle("/{review_id}", utils.AuthMiddleware(http.HandlerFunc(controller.EditReview))).Methods("OPTIONS", "PUT")
	reviewRouter.Handle("/{review_id}", utils.AuthMiddleware(http.HandlerFunc(controller.DeleteReview))).Methods("OPTIONS", "DELETE")

	favoriteRouter := router.PathPrefix("/favorites").Subrouter()                                                        // Ganti "/favorite" menjadi "/favorites"
	favoriteRouter.Handle("/{anime_id}", utils.AuthMiddleware(http.HandlerFunc(controller.AddFavorite))).Methods("POST") // Menambahkan favorit
	favoriteRouter.Handle("/", utils.AuthMiddleware(http.HandlerFunc(controller.GetFavorites))).Methods("GET")           // Mengambil daftar favorit
	favoriteRouter.Handle("/{id}", utils.AuthMiddleware(http.HandlerFunc(controller.RemoveFavorite))).Methods("DELETE")  // Menghapus favorit berdasarkan ID

	return router
}
