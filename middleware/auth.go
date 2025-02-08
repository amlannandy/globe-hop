package middleware

import (
	"context"
	"globe-hop/config"
	"globe-hop/models"
	"globe-hop/utils"
	"net/http"
	"strings"
)

// AuthMiddleware is a middleware function that authenticates requests using JWT tokens
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the Authorization header from the request
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.SendErrorResponse(w, &utils.ErrorResponseBody{
				Status: http.StatusUnauthorized,
				Error:  "Authorization header is required",
			})
			return
		}

		// Split the header to get the token part
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.SendErrorResponse(w, &utils.ErrorResponseBody{
				Status: http.StatusUnauthorized,
				Error:  "Invalid authorization header format",
			})
			return
		}

		// Get the token
		token := headerParts[1]

		// Decode token and fetch user ID
		userId, err := config.DecodeJWTToken(token)
		if err != nil {
			utils.SendErrorResponse(w, &utils.ErrorResponseBody{
				Status: http.StatusUnauthorized,
				Error:  "Invalid token",
			})
			return
		}

		// Retrieve the user from the database using the extracted user ID
		var user *models.User
		err = config.DB.Where("id = ?", userId).First(&user).Error
		if err != nil {
			utils.SendErrorResponse(w, &utils.ErrorResponseBody{
				Status: http.StatusUnauthorized,
				Error:  "User not found",
			})
			return
		}
		// Add userId to request context
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
