package authz

import (
	"context"
	"net/http"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	rolesKey  contextKey = "roles"
	scopesKey contextKey = "scopes"
)

// Role constants
const (
	RoleAdmin   = "admin"
	RoleTeacher = "teacher"
	RoleLearner = "learner"
)

func WithAuth(ctx context.Context, userID string, roles, scopes []string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, rolesKey, roles)
	ctx = context.WithValue(ctx, scopesKey, scopes)
	return ctx
}

func RolesFromContext(ctx context.Context) []string {
	roles, _ := ctx.Value(rolesKey).([]string)
	return roles
}

func ScopesFromContext(ctx context.Context) []string {
	scopes, _ := ctx.Value(scopesKey).([]string)
	return scopes
}

func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles := RolesFromContext(r.Context())
			if len(roles) == 0 {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			hasRole := false
			for _, role := range roles {
				if role == requiredRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireScope(requiredScope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scopes := ScopesFromContext(r.Context())
			if len(scopes) == 0 {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			hasScope := false
			for _, scope := range scopes {
				if scope == requiredScope {
					hasScope = true
					break
				}
			}

			if !hasScope {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func InferScopesFromRoles(roles []string) []string {
	hasWrite := false
	for _, role := range roles {
		switch role {
		case RoleAdmin, RoleTeacher, RoleLearner:
			hasWrite = true
		}
	}
	if hasWrite {
		return []string{"read", "write"}
	}
	return []string{"read"}
}
