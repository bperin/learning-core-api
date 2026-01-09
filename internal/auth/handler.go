package auth

import (
	"net/http"

	authdto "learning-core-api/internal/auth/dto"
	"learning-core-api/internal/http/render"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	authService *Service
}

func NewHandler(authService *Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	r.Post("/oauth/token", h.Token)
}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {}

// Token godoc
// @Summary OAuth token endpoint
// @Description Issue or refresh OAuth2 tokens using supported grant types.
// @Description Supported grant types:
// @Description - password: exchange email + password for tokens
// @Description - refresh_token: exchange refresh_token for new access token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authdto.TokenRequest true "OAuth Token Request"
// @Success 200 {object} authdto.TokenResponse "OAuth2 token response"
// @Failure 400 {string} string "invalid request or unsupported grant_type"
// @Failure 401 {string} string "invalid credentials"
// @Router /oauth/token [post]
func (h *Handler) Token(w http.ResponseWriter, r *http.Request) {
	var req authdto.TokenRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	switch req.GrantType {
	case "password":
		if req.Email == nil || req.Password == nil || *req.Email == "" || *req.Password == "" {
			render.Error(w, http.StatusBadRequest, "email and password are required for grant_type=password")
			return
		}

		tokens, _, err := h.authService.LoginWithEmail(r.Context(), *req.Email, *req.Password)
		if err != nil {
			render.Error(w, http.StatusUnauthorized, err.Error())
			return
		}

		render.JSON(w, http.StatusOK, tokens)
		return

	case "refresh_token":
		if req.RefreshToken == nil || *req.RefreshToken == "" {
			render.Error(w, http.StatusBadRequest, "refresh_token is required for grant_type=refresh_token")
			return
		}

		tokens, err := h.authService.RefreshToken(r.Context(), *req.RefreshToken)
		if err != nil {
			render.Error(w, http.StatusUnauthorized, err.Error())
			return
		}

		render.JSON(w, http.StatusOK, tokens)
		return

	default:
		render.Error(w, http.StatusBadRequest, "unsupported grant_type")
		return
	}
}
