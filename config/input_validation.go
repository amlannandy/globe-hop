package config

import "github.com/go-playground/validator/v10"

// Helper function to format validation errors
func FormatValidationError(err error) string {
	var errorMessages []string
	for _, err := range err.(validator.ValidationErrors) {
		errorMessages = append(errorMessages, err.Field()+" is "+err.Tag())
	}
	return "Validation failed: " + stringJoin(errorMessages, ", ")
}

// Helper function to join strings with a separator
func stringJoin(strs []string, sep string) string {
	result := ""
	for i, str := range strs {
		if i > 0 {
			result += sep
		}
		result += str
	}
	return result
}
