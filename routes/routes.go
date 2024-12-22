package routes

import (
	"NYANIMEBACKEND/controller"
	"NYANIMEBACKEND/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Redirect root to the frontend
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://nya-nime.github.io/Nyanime/", http.StatusFound)
	}).Methods("GET")

	// User Routes
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/register", controller.Register).Methods("OPTIONS", "POST")
	userRouter.HandleFunc("/login", controller.Login).Methods("POST")
	userRouter.Handle("/profile", utils.AuthMiddleware(http.HandlerFunc(controller.GetUserProfile))).Methods("GET")
	userRouter.Handle("/logout", utils.AuthMiddleware(http.HandlerFunc(controller.Logout))).Methods("POST")

	// Anime Routes (Admin Privileges)
	animeRouter := router.PathPrefix("/anime").Subrouter()
	animeRouter.HandleFunc("/", controller.GetAllAnime).Methods("GET")                                                                   // Semua user bisa lihat
	animeRouter.Handle("/", utils.AuthMiddleware(utils.AdminMiddleware(http.HandlerFunc(controller.CreateAnime)))).Methods("POST")       // Menambahkan anime
	animeRouter.Handle("/{id}", utils.AuthMiddleware(utils.AdminMiddleware(http.HandlerFunc(controller.EditAnime)))).Methods("PUT")      // Mengedit anime
	animeRouter.Handle("/{id}", utils.AuthMiddleware(utils.AdminMiddleware(http.HandlerFunc(controller.DeleteAnime)))).Methods("DELETE") // Menghapus anime

	// Review Routes
	reviewRouter := router.PathPrefix("/review").Subrouter()
	reviewRouter.Handle("/anime/{anime_id}", utils.AuthMiddleware(http.HandlerFunc(controller.AddReview))).Methods("POST")
	reviewRouter.Handle("/{id}", utils.AuthMiddleware(http.HandlerFunc(controller.EditReview))).Methods("PUT")
	reviewRouter.Handle("/{id}", utils.AuthMiddleware(http.HandlerFunc(controller.DeleteReview))).Methods("DELETE")

	return router
}
