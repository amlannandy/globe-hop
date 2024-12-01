package routes

import "github.com/gorilla/mux"

func InitializeRouter() *mux.Router {
	router := mux.NewRouter()

	InitializeUsersRoutes(router)

	return router
}
