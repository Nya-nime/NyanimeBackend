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

	// Review Routes
	reviewRouter := router.PathPrefix("/review").Subrouter()
	reviewRouter.Handle("/anime/{anime_id}", utils.AuthMiddleware(http.HandlerFunc(controller.AddReview))).Methods("POST")
	reviewRouter.Handle("/{id}", utils.AuthMiddleware(http.HandlerFunc(controller.EditReview))).Methods("PUT")
	reviewRouter.Handle("/{id}", utils.AuthMiddleware(http.HandlerFunc(controller.DeleteReview))).Methods("DELETE")

	return router
}
