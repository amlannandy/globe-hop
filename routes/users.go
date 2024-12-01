package routes

import (
	"globe-hop/controllers"

	"github.com/gorilla/mux"
)

func InitializeUsersRoutes(r *mux.Router) {
	r.HandleFunc("/users/login", controllers.Login).Methods("POST")
	r.HandleFunc("/users/register", controllers.Register).Methods("POST")
}
