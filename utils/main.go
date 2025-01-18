package utils

import (
	"encoding/json"
	"net/http"
)

type SuccessResponseBody struct {
	Status  int
	Message string
	Data    any
}

type ErrorResponseBody struct {
	Status int
	Error  string
}

func SendSuccessResponse(w http.ResponseWriter, data *SuccessResponseBody) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(data.Status)
	json.NewEncoder(w).Encode(map[string]any{
		"message": data.Message,
		"data":    data.Data,
	})
}

func SendErrorResponse(w http.ResponseWriter, data *ErrorResponseBody) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(data.Status)
	json.NewEncoder(w).Encode(map[string]any{
		"error": data.Error,
	})
}
