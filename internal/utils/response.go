package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate is a validator instance
var Validate = validator.New()

// RespondWithError responds with an error
func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

// RespondWithJSON responds with JSON
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal Server Error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// ValidationErrorMessage returns a formatted validation error message
func ValidationErrorMessage(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range validationErrors {
			field := e.Field()
			switch e.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("%s is required", field))
			case "email":
				messages = append(messages, fmt.Sprintf("%s must be a valid email", field))
			case "min":
				messages = append(messages, fmt.Sprintf("%s must be at least %s characters", field, e.Param()))
			case "max":
				messages = append(messages, fmt.Sprintf("%s must be at most %s characters", field, e.Param()))
			case "oneof":
				messages = append(messages, fmt.Sprintf("%s must be one of: %s", field, e.Param()))
			default:
				messages = append(messages, fmt.Sprintf("%s is invalid", field))
			}
		}
		return strings.Join(messages, ", ")
	}
	return err.Error()
}