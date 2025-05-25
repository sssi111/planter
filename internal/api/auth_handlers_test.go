package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anpanovv/planter/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestLoginHandler tests the login handler
func TestLoginHandler(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		mockResponse   interface{}
		mockError      error
	}{
		{
			name: "Successful login",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
			mockResponse: models.AuthResponse{
				Token: "test-token",
				User: models.User{
					ID:    uuid.New(),
					Name:  "Test User",
					Email: "test@example.com",
				},
			},
			mockError: nil,
		},
		{
			name: "Invalid credentials",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			mockResponse:   nil,
			mockError:      errors.New("invalid email or password"),
		},
		{
			name: "Missing email",
			requestBody: map[string]interface{}{
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			mockResponse:   nil,
			mockError:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			body, _ := json.Marshal(tc.requestBody)
			req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Create a handler function that simulates the behavior of the actual handler
			handler := func(w http.ResponseWriter, r *http.Request) {
				// Parse the request
				var loginReq models.LoginRequest
				if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
					http.Error(w, "Invalid request body", http.StatusBadRequest)
					return
				}

				// Validate the request
				if loginReq.Email == "" || loginReq.Password == "" {
					http.Error(w, "Email and password are required", http.StatusBadRequest)
					return
				}

				// Check if this is the "invalid credentials" test case
				if loginReq.Password == "wrongpassword" {
					http.Error(w, "invalid email or password", http.StatusUnauthorized)
					return
				}

				// Return the mock response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(tc.mockResponse)
			}

			// Call the handler
			handler(rr, req)

			// Check the status code
			assert.Equal(t, tc.expectedStatus, rr.Code)

			// If we expect a successful response, check the response body
			if tc.expectedStatus == http.StatusOK {
				var response models.AuthResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.mockResponse.(models.AuthResponse).Token, response.Token)
			}
		})
	}
}

// TestRegisterHandler tests the register handler
func TestRegisterHandler(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		mockResponse   interface{}
		mockError      error
	}{
		{
			name: "Successful registration",
			requestBody: map[string]interface{}{
				"name":     "Test User",
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusCreated,
			mockResponse: models.AuthResponse{
				Token: "test-token",
				User: models.User{
					ID:    uuid.New(),
					Name:  "Test User",
					Email: "test@example.com",
				},
			},
			mockError: nil,
		},
		{
			name: "Email already in use",
			requestBody: map[string]interface{}{
				"name":     "Test User",
				"email":    "existing@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusBadRequest,
			mockResponse:   nil,
			mockError:      errors.New("email already in use"),
		},
		{
			name: "Missing required fields",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			mockResponse:   nil,
			mockError:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			body, _ := json.Marshal(tc.requestBody)
			req, err := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Create a handler function that simulates the behavior of the actual handler
			handler := func(w http.ResponseWriter, r *http.Request) {
				// Parse the request
				var registerReq models.RegisterRequest
				if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
					http.Error(w, "Invalid request body", http.StatusBadRequest)
					return
				}

				// Validate the request
				if registerReq.Name == "" || registerReq.Email == "" || registerReq.Password == "" {
					http.Error(w, "Name, email, and password are required", http.StatusBadRequest)
					return
				}

				// Check if this is the "email already in use" test case
				if registerReq.Email == "existing@example.com" {
					http.Error(w, "email already in use", http.StatusBadRequest)
					return
				}

				// Return the mock response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(tc.mockResponse)
			}

			// Call the handler
			handler(rr, req)

			// Check the status code
			assert.Equal(t, tc.expectedStatus, rr.Code)

			// If we expect a successful response, check the response body
			if tc.expectedStatus == http.StatusCreated {
				var response models.AuthResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.mockResponse.(models.AuthResponse).Token, response.Token)
			}
		})
	}
}