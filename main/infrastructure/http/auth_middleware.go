package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const (
	jwksURL        = "http://localhost:3000/api/auth/jwks"
	jwksCacheTTL   = 1 * time.Hour
	contextUserKey = "user"
)

var (
	jwksCache    jwk.Set
	jwksCacheExp time.Time
)

// UserClaims represents the authenticated user information extracted from the JWT
type UserClaims struct {
	UserID string
	Email  string
}

// fetchJWKS fetches the JWKS from the auth server
func fetchJWKS() (jwk.Set, error) {
	// Check cache first
	if jwksCache != nil && time.Now().Before(jwksCacheExp) {
		return jwksCache, nil
	}

	// Fetch JWKS from the auth server
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: status code %d", resp.StatusCode)
	}

	// Parse the JWKS
	set, err := jwk.ParseReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %w", err)
	}

	// Cache the JWKS
	jwksCache = set
	jwksCacheExp = time.Now().Add(jwksCacheTTL)

	return set, nil
}

// extractBearerToken extracts the Bearer token from the Authorization header
func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	// Check if it's a Bearer token
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	return parts[1], nil
}

// verifyToken verifies the JWT token using the JWKS
func verifyToken(tokenString string) (jwt.Token, error) {
	// Fetch JWKS
	keySet, err := fetchJWKS()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}

	// Parse and verify the token
	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	return token, nil
}

// extractUserClaims extracts user information from the JWT token
func extractUserClaims(token jwt.Token) (*UserClaims, error) {
	// Get the subject (user ID)
	subject := token.Subject()
	if subject == "" {
		return nil, fmt.Errorf("token missing subject claim")
	}

	// Get the email claim
	email, ok := token.Get("email")
	if !ok {
		return nil, fmt.Errorf("token missing email claim")
	}

	emailStr, ok := email.(string)
	if !ok {
		return nil, fmt.Errorf("email claim is not a string")
	}

	return &UserClaims{
		UserID: subject,
		Email:  emailStr,
	}, nil
}

// GetUserFromContext retrieves the authenticated user from the request context
func GetUserFromContext(ctx context.Context) (*UserClaims, bool) {
	user, ok := ctx.Value(contextUserKey).(*UserClaims)
	return user, ok
}

// authMiddleware validates JWT tokens and adds user information to the request context
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the Bearer token
		tokenString, err := extractBearerToken(r)
		if err != nil {
			log.Printf("Authentication failed: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Unauthorized: " + err.Error(),
			})
			return
		}

		// Verify the token
		token, err := verifyToken(tokenString)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Unauthorized: invalid or expired token",
			})
			return
		}

		// Extract user claims
		userClaims, err := extractUserClaims(token)
		if err != nil {
			log.Printf("Failed to extract user claims: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Unauthorized: invalid token claims",
			})
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), contextUserKey, userClaims)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
