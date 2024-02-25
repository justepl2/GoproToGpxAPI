package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

func ValidateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Get the JSON Web Key Set (JWKS) from AWS Cognito

		set, err := jwk.Fetch(context.Background(), "https://cognito-idp."+os.Getenv("AWS_REGION")+".amazonaws.com/"+os.Getenv("AWS_COGNITO_USER_POOL_ID")+"/.well-known/jwks.json")
		if err != nil {
			http.Error(w, "Failed to fetch JWKS", http.StatusInternalServerError)
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse([]byte(parts[1]), jwt.WithKeySet(set))
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID, exists := token.Get("sub")
		if !exists {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userID)

		// The token is valid, pass the request to the next middleware
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LogoutUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "token", parts[1])

		// The token is valid, pass the request to the next middleware
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
