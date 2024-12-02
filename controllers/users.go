package controllers

import (
	"encoding/json"
	"globe-hop/config"
	"globe-hop/models"
	"globe-hop/utils"
	"net/http"

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
			Error:  "Invalid request",
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
			Error:  "Error processing password",
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
			Error:  "Email is already registered",
		})
		return
	}

	// Save user to database
	err = config.DB.Create(&user).Error
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusInternalServerError,
			Error:  "Error registering user",
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

}
