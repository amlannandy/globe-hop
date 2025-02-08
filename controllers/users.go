package controllers

import (
	"encoding/json"
	"globe-hop/config"
	"globe-hop/models"
	"globe-hop/types"
	"globe-hop/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

// Register registers a new user with the provided email and password.
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

// Login logs in a user with the provided email and password.
func Login(w http.ResponseWriter, r *http.Request) {
	var body types.LoginRequestBody

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusBadRequest,
			Error:  "Invalid request",
		})
		return
	}

	// Validate body
	err = validate.Struct(body)
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusBadRequest,
			Error:  config.FormatValidationError(err),
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

	// Check if the password is correct
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

// GetCurrentUser retrieves the current user based on the JWT token provided in the request header.
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {

	// Send success response with the current user data
	utils.SendSuccessResponse(w, &utils.SuccessResponseBody{
		Status:  http.StatusAccepted,
		Message: "Current user retrieved.",
		Data:    r.Context().Value("user"),
	})
}

// DeleteUser deletes the current user.
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	// Get password from request body
	var body types.DeleteUserRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusBadRequest,
			Error:  "Invalid request body",
		})
		return
	}

	// Validate body
	err := validate.Struct(body)
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusBadRequest,
			Error:  config.FormatValidationError(err),
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusUnauthorized,
			Error:  "Incorrect password",
		})
		return
	}

	// Delete user from database
	err = config.DB.Unscoped().Delete(user).Error
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusInternalServerError,
			Error:  "Error deleting user.",
		})
		return
	}

	// Return success response
	utils.SendSuccessResponse(w, &utils.SuccessResponseBody{
		Status:  http.StatusOK,
		Message: "User deleted successfully.",
	})
}

// UpdatePassword updates the password of the current user.
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*models.User)

	// Get password from request body
	var body types.UpdatePasswordRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusBadRequest,
			Error:  "Invalid request body",
		})
		return
	}

	// Validate body
	err := validate.Struct(body)
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusBadRequest,
			Error:  config.FormatValidationError(err),
		})
		return
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.OldPassword)); err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusUnauthorized,
			Error:  "Incorrect old password",
		})
		return
	}

	// Check if old and new passwords are same
	if body.OldPassword == body.NewPassword {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusBadRequest,
			Error:  "New password must be different from old password",
		})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusInternalServerError,
			Error:  "Error updating password",
		})
		return
	}

	// Update user password
	user.Password = string(hashedPassword)
	err = config.DB.Save(user).Error
	if err != nil {
		utils.SendErrorResponse(w, &utils.ErrorResponseBody{
			Status: http.StatusInternalServerError,
			Error:  "Error updating password",
		})
		return
	}

	// Return success response
	utils.SendSuccessResponse(w, &utils.SuccessResponseBody{
		Status:  http.StatusOK,
		Message: "Password updated successfully.",
	})
}
