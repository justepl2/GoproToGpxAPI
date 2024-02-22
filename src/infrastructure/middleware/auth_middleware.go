package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/justepl2/gopro_to_gpx_api/domain"
	"github.com/justepl2/gopro_to_gpx_api/tools"
)

func UserAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Récupérer le token de l'en-tête Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Le token doit être préfixé par "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			tools.FormatResponseBody(w, http.StatusUnauthorized, "Invalid Authorization header")

			return
		}

		// Vérifier le token
		segments := strings.Split(parts[1], ".")
		if len(segments) != 3 {
			tools.FormatResponseBody(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			tools.FormatResponseBody(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			tools.FormatResponseBody(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			tools.FormatResponseBody(w, http.StatusUnauthorized, "Invalid role claim")
			return
		}

		// Vérifier le rôle
		if role != string(domain.RoleUser) && role != string(domain.RoleAdmin) {
			tools.FormatResponseBody(w, http.StatusForbidden, "Forbidden")
			return
		}

		// L'utilisateur est authentifié, passer à la prochaine requête
		next.ServeHTTP(w, r)
	})
}

func AdminAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Récupérer le token de l'en-tête Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Le token doit être préfixé par "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			tools.FormatResponseBody(w, http.StatusUnauthorized, "Invalid Authorization header")

			return
		}

		// Vérifier le token
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("your-secret"), nil
		})
		if err != nil || !token.Valid {
			tools.FormatResponseBody(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			tools.FormatResponseBody(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			tools.FormatResponseBody(w, http.StatusUnauthorized, "Invalid role claim")
			return
		}

		// Vérifier le rôle
		if role != string(domain.RoleAdmin) {
			tools.FormatResponseBody(w, http.StatusForbidden, "Forbidden")
			return
		}

		// L'utilisateur est authentifié, passer à la prochaine requête
		next.ServeHTTP(w, r)
	})
}
