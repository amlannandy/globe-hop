package controllers

import (
	"encoding/json"
	"globe-hop/config"
	"globe-hop/models"
	"globe-hop/types"
	"globe-hop/utils"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusBadRequest,
			Error:  "Invalid request.",
		})
		return
	}

	// Validate user input
	err = validate.Struct(user)
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusBadRequest,
			Error:  config.FormatValidationError(err),
		})
		return
	}

	// Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusInternalServerError,
			Error:  "Error processing password.",
		})
		return
	}
	user.Password = string(hashedPassword)

	// Check if the email is already registered
	var existingUser *models.User
	err = config.DB.Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusConflict,
			Error:  "Email is already registered.",
		})
		return
	}

	// Save user to database
	err = config.DB.Create(&user).Error
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusInternalServerError,
			Error:  "Error registering user.",
		})
		return
	}

	// Generate JWT token
	token, err := config.GenerateJWTToken(&user)
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusInternalServerError,
			Error:  "Error generating token.",
		})
		return
	}

	// Return success response with token
	utils.SendSuccessResponse(w, &utils.SuccessResponseBody{
		Status:  http.StatusAccepted,
		Message: "User registered successfully.",
		Data:    token,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var body types.LoginRequestBody

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusBadRequest,
			Error:  "Invalid request",
		})
		return
	}

	// Check if user exists or not
	var user *models.User
	err = config.DB.Where("email = ?", body.Email).First(&user).Error
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusNotFound,
			Error:  "Account with this email does not exist",
		})
		return
	}

	authResult := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if authResult != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusUnauthorized,
			Error:  "Incorrect password",
		})
		return
	}

	// Generate JWT token
	token, err := config.GenerateJWTToken(user)
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusInternalServerError,
			Error:  "Error generating token.",
		})
		return
	}

	// Return success response with token
	utils.SendSuccessResponse(w, &utils.SuccessResponseBody{
		Status:  http.StatusAccepted,
		Message: "User logged in successfully.",
		Data:    token,
	})
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token := strings.Split(authHeader, " ")[1]

	// Decode token and fetch user
	userId, err := config.DecodeJWTToken(token)
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusUnauthorized,
			Error:  "Invalid token.",
		})
		return
	}

	var user *models.User
	err = config.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusNotFound,
			Error:  "User not found.",
		})
		return
	}

	// Return success response with token
	utils.SendSuccessResponse(w, &utils.SuccessResponseBody{
		Status:  http.StatusAccepted,
		Message: "Current user retrieved.",
		Data:    *user,
	})
}
