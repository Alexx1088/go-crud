package middleware

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
)

type contextKey string

const (
	userIDKey contextKey = "user_id" // Key to store user_id in the context
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrMissingAuth  = errors.New("missing Authorization header")
)

// AuthMiddleware validates JWT tokens and ensures requests are authenticated
func AuthMiddleware(secretKey []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Middleware: Starting token validation...")
			// Step 1: Extract the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Println("Middleware: Missing Authorization header")
				http.Error(w, ErrMissingAuth.Error(), http.StatusUnauthorized)
				return
			}
			// Step 2: Validate the header format (Bearer <token>)
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader { // No "Bearer " prefix found
				log.Println("Middleware: Invalid Authorization header format")
				http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			// Step 3: Parse and verify the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Ensure the signing method is HMAC
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Println("Middleware: Invalid signing method")
					return nil, ErrInvalidToken
				}
				return secretKey, nil
			})

			if err != nil || !token.Valid {
				log.Printf("Middleware: Token validation failed: %v", err)
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}
			// Step 4: Extract claims and validate them
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Println("Middleware: Invalid token claims")
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			// Extract user_id from claims
			userID, ok := claims["user_id"].(float64) // JWT stores numbers as float64
			if !ok {
				log.Println("Middleware: Missing or invalid user_id claim")
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}

			log.Printf("Middleware: Token validated successfully for user_id: %d", int(userID))

			// Step 5: Add user_id to the request context
			ctx := context.WithValue(r.Context(), userIDKey, int(userID))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
