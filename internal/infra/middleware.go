package infra

import (
	"fmt"
	"net/http"
	"strings"

	"learning-core-api/internal/http/authz"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid claims", http.StatusUnauthorized)
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}

			// Extract roles from claims
			var roles []string
			if rClaim, ok := claims["roles"]; ok {
				switch v := rClaim.(type) {
				case []interface{}:
					for _, r := range v {
						if s, ok := r.(string); ok {
							roles = append(roles, s)
						}
					}
				case string:
					roles = strings.Split(v, ",")
				}
			}

			ctx := authz.WithAuth(r.Context(), userID, roles, extractScopes(claims, roles))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractScopes(claims jwt.MapClaims, roles []string) []string {
	var scopes []string
	if sClaim, ok := claims["scopes"]; ok {
		switch v := sClaim.(type) {
		case []interface{}:
			for _, s := range v {
				if str, ok := s.(string); ok {
					scopes = append(scopes, str)
				}
			}
		case string:
			scopes = append(scopes, strings.Split(v, " ")...)
		}
	}
	if len(scopes) == 0 {
		if sClaim, ok := claims["scope"]; ok {
			if v, ok := sClaim.(string); ok {
				scopes = append(scopes, strings.Split(v, " ")...)
			}
		}
	}
	if len(scopes) == 0 {
		return authz.InferScopesFromRoles(roles)
	}
	return scopes
}
