package api

import (
	"net/http"
	"user-management-system/pkg/handlers"
	"user-management-system/pkg/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Apply CORS middleware
	router.Use(middleware.CORS)

	router.HandleFunc("/signup", handlers.SignUp).Methods(http.MethodPost)
	router.HandleFunc("/signin", handlers.SignIn).Methods(http.MethodPost)
	// Add other routes here...

	return router
}
