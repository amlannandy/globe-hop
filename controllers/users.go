package controllers

import (
	"encoding/json"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted);
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Login route",
		"status": 200,
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted);
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Register route",
		"status": 200,
	})
}