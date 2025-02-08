package routes

import (
	"globe-hop/controllers"
	"globe-hop/middleware"

	"github.com/gorilla/mux"
)

func InitializeUsersRoutes(r *mux.Router) {
	r.HandleFunc("/users/login", controllers.Login).Methods("POST")
	r.HandleFunc("/users/register", controllers.Register).Methods("POST")
	r.HandleFunc("/users/current", middleware.AuthMiddleware(controllers.GetCurrentUser)).Methods("GET")
	r.HandleFunc("/users/delete", middleware.AuthMiddleware(controllers.DeleteUser)).Methods("DELETE")
}
