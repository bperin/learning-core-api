package infra

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"

	"learning-core-api/internal/auth"
	"learning-core-api/internal/domain/attempts"
	"learning-core-api/internal/domain/chunking_configs"
	"learning-core-api/internal/domain/documents"
	"learning-core-api/internal/domain/evals"
	"learning-core-api/internal/domain/prompt_templates"
	"learning-core-api/internal/domain/reviews"
	"learning-core-api/internal/domain/schema_templates"
	"learning-core-api/internal/domain/system_instructions"
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timestamp := time.Now().Format(time.RFC3339)
		
		log.Printf("\n[%s] [API Request] %s %s", timestamp, r.Method, r.URL.Path)
		log.Printf("[API Headers] %v", r.Header)
		
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenPreview := authHeader
			if len(tokenPreview) > 50 {
				tokenPreview = tokenPreview[:50] + "..."
			}
			log.Printf("[API Auth] ✓ Authorization header present: %s", tokenPreview)
		} else {
			log.Printf("[API Auth] ✗ NO Authorization header!")
		}
		
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			if len(bodyBytes) > 0 {
				log.Printf("[API Request Body] %s", string(bodyBytes))
			}
		}
		
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		next.ServeHTTP(wrappedWriter, r)
		
		log.Printf("[API Response] %d %s %s", wrappedWriter.statusCode, r.Method, r.URL.Path)
		log.Printf("[API Response Headers] %v", wrappedWriter.Header())
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func NewRouter(deps RouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(loggingMiddleware)
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

	// Schema management handlers
	promptTemplatesRepo := prompt_templates.NewRepository(deps.Queries)
	promptTemplatesService := prompt_templates.NewService(promptTemplatesRepo)
	promptTemplatesHandler := prompt_templates.NewHandler(promptTemplatesService)

	schemaTemplatesRepo := schema_templates.NewRepository(deps.Queries)
	schemaTemplatesService := schema_templates.NewService(schemaTemplatesRepo)
	schemaTemplatesHandler := schema_templates.NewHandler(schemaTemplatesService)

	chunkingConfigsRepo := chunking_configs.NewRepository(deps.Queries)
	chunkingConfigsService := chunking_configs.NewService(chunkingConfigsRepo)
	chunkingConfigsHandler := chunking_configs.NewHandler(chunkingConfigsService)

	systemInstructionsRepo := system_instructions.NewRepository(deps.Queries)
	systemInstructionsService := system_instructions.NewService(systemInstructionsRepo)
	systemInstructionsHandler := system_instructions.NewHandler(systemInstructionsService)

	authHandler.RegisterPublicRoutes(r)
	usersHandler.RegisterPublicRoutes(r)

	registerRoleRoutes(r, deps.JWTSecret, authHandler)
	registerRoleRoutes(r, deps.JWTSecret, usersHandler)
	registerRoleRoutes(r, deps.JWTSecret, documentsHandler)
	registerRoleRoutes(r, deps.JWTSecret, textbooksHandler)
	registerRoleRoutes(r, deps.JWTSecret, evalsHandler)
	registerRoleRoutes(r, deps.JWTSecret, reviewsHandler)
	registerRoleRoutes(r, deps.JWTSecret, attemptsHandler)
	registerRoleRoutes(r, deps.JWTSecret, promptTemplatesHandler)
	registerRoleRoutes(r, deps.JWTSecret, schemaTemplatesHandler)
	registerRoleRoutes(r, deps.JWTSecret, chunkingConfigsHandler)
	registerRoleRoutes(r, deps.JWTSecret, systemInstructionsHandler)

	return r
}

func registerRoleRoutes(r chi.Router, secret string, registrar RoleRouteRegistrar) {
	registerProtectedRoleRoutes(r, secret, authz.RoleLearner, registrar.RegisterLearnerRoutes)
	registerProtectedRoleRoutes(r, secret, authz.RoleTeacher, registrar.RegisterTeacherRoutes)
	registerProtectedRoleRoutes(r, secret, authz.RoleAdmin, registrar.RegisterAdminRoutes)
}

func registerProtectedRoleRoutes(r chi.Router, secret, role string, register func(chi.Router)) {
	r.Group(func(r chi.Router) {
		r.Use(JWTMiddleware(secret))
		r.Use(authz.RequireRole(role))
		register(r)
	})
}
