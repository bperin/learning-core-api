package infra

import (
	"net/http"

	"learning-core-api/internal/auth"
	"learning-core-api/internal/domain/attempts"
	"learning-core-api/internal/domain/documents"
	"learning-core-api/internal/domain/evals"
	"learning-core-api/internal/domain/reviews"
	"learning-core-api/internal/domain/textbooks"
	"learning-core-api/internal/domain/users"
	"learning-core-api/internal/gcp"
	"learning-core-api/internal/http/authz"
	"learning-core-api/internal/persistance/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterDeps struct {
	JWTSecret    string
	Queries      *store.Queries
	GoogleAPIKey string
	GCSService   *gcp.GCSService
}

type RoleRouteRegistrar interface {
	RegisterPublicRoutes(r chi.Router)
	RegisterAdminRoutes(r chi.Router)
	RegisterTeacherRoutes(r chi.Router)
	RegisterLearnerRoutes(r chi.Router)
}

func NewRouter(deps RouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	authRepo := auth.NewRepository(deps.Queries)
	authService := auth.NewService(deps.JWTSecret, authRepo)
	authHandler := auth.NewHandler(authService)

	usersRepo := users.NewRepository(deps.Queries)
	usersService := users.NewService(usersRepo)
	usersHandler := users.NewHandler(usersService)

	documentsService := documents.NewService(documents.NewRepository(deps.Queries))
	documentsHandler := documents.NewHandler(documentsService, deps.GCSService)

	textbooksRepo := textbooks.NewRepository()
	textbooksService := textbooks.NewService(textbooksRepo)
	textbooksHandler := textbooks.NewHandler(textbooksService)

	evalsHandler := evals.NewHandler()
	reviewsHandler := reviews.NewHandler()
	attemptsHandler := attempts.NewHandler()

	authHandler.RegisterPublicRoutes(r)
	usersHandler.RegisterPublicRoutes(r)

	registerRoleRoutes(r, deps.JWTSecret, authHandler)
	registerRoleRoutes(r, deps.JWTSecret, usersHandler)
	registerRoleRoutes(r, deps.JWTSecret, documentsHandler)
	registerRoleRoutes(r, deps.JWTSecret, textbooksHandler)
	registerRoleRoutes(r, deps.JWTSecret, evalsHandler)
	registerRoleRoutes(r, deps.JWTSecret, reviewsHandler)
	registerRoleRoutes(r, deps.JWTSecret, attemptsHandler)

	return r
}

func registerRoleRoutes(r chi.Router, secret string, registrar RoleRouteRegistrar) {
	registerProtectedRoleRoutes(r, secret, authz.RoleAdmin, registrar.RegisterAdminRoutes)
	registerProtectedRoleRoutes(r, secret, authz.RoleTeacher, registrar.RegisterTeacherRoutes)
	registerProtectedRoleRoutes(r, secret, authz.RoleLearner, registrar.RegisterLearnerRoutes)
}

func registerProtectedRoleRoutes(r chi.Router, secret, role string, register func(chi.Router)) {
	r.Group(func(r chi.Router) {
		r.Use(JWTMiddleware(secret))
		r.Use(authz.RequireRole(role))
		register(r)
	})
}
