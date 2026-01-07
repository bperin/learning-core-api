package infra

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"learning-core-api/internal/documents"
	"learning-core-api/internal/filesearch"
	"learning-core-api/internal/store"
	"learning-core-api/internal/subjects"
)

type RouterDeps struct {
	JWTSecret    string
	Queries      *store.Queries
	GoogleAPIKey string
}

func NewRouter(deps RouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(JWTMiddleware(deps.JWTSecret))

		r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(UserIDKey)
			w.Write([]byte(fmt.Sprintf("User ID: %v", userID)))
		})

		if deps.Queries == nil {
			return
		}

		subjectRepo := subjects.NewRepository(deps.Queries)
		subjectService := subjects.NewService(subjectRepo)
		subjectHandler := subjects.NewHandler(subjectService)

		documentRepo := documents.NewRepository(deps.Queries)
		documentService := documents.NewService(documentRepo)
		documentHandler := documents.NewHandler(documentService)

		fileSearchRepo := filesearch.NewRepository(deps.Queries)
		fileSearchService := filesearch.NewService(deps.GoogleAPIKey, fileSearchRepo, subjectRepo, documentRepo)
		fileSearchHandler := filesearch.NewHandler(fileSearchService)

		r.Route("/subjects", subjectHandler.RegisterRoutes)
		r.Route("/documents", documentHandler.RegisterRoutes)
		r.Route("/file-search", fileSearchHandler.RegisterRoutes)
	})

	return r
}
