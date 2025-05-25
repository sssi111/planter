package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// contextKey is a custom type for context keys
type contextKey string

// UserIDKey is the key for user ID in the request context
const UserIDKey contextKey = "userID"

// JWTClaims represents the claims in a JWT
type JWTClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

// Auth is the authentication middleware
type Auth struct {
	jwtSecret string
}

// NewAuth creates a new Auth middleware
func NewAuth(jwtSecret string) *Auth {
	return &Auth{
		jwtSecret: jwtSecret,
	}
}

// Middleware authenticates the request
func (a *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
			return
		}

		// Parse the token
		token := parts[1]
		claims, err := a.parseToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add the user ID to the request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// parseToken parses and validates a JWT token
func (a *Auth) parseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateToken generates a JWT token for a user
func (a *Auth) GenerateToken(userID uuid.UUID, duration time.Duration) (string, error) {
	// Create the claims
	claims := &JWTClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	return token.SignedString([]byte(a.jwtSecret))
}

// GetUserID gets the user ID from the request context
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userIDStr, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID in context")
	}

	return userID, nil
}

// RequireAuth is a middleware that requires authentication
func (a *Auth) RequireAuth(next http.Handler) http.Handler {
	return a.Middleware(next)
}

// OptionalAuth is a middleware that makes authentication optional
func (a *Auth) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No authentication, just continue
			next.ServeHTTP(w, r)
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Invalid format, just continue without authentication
			next.ServeHTTP(w, r)
			return
		}

		// Parse the token
		token := parts[1]
		claims, err := a.parseToken(token)
		if err != nil {
			// Invalid token, just continue without authentication
			next.ServeHTTP(w, r)
			return
		}

		// Add the user ID to the request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}