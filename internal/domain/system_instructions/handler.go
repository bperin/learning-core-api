package system_instructions

import (
	"net/http"

	"learning-core-api/internal/http/authz"
	"learning-core-api/internal/http/render"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterPublicRoutes(r chi.Router) {
	// No public routes for system instructions
}

func (h *Handler) RegisterAdminRoutes(r chi.Router) {
	r.With(authz.RequireScope("read")).Get("/system-instructions", h.ListAll)
	r.With(authz.RequireScope("read")).Get("/system-instructions/{id}", h.GetByID)
	r.With(authz.RequireScope("read")).Get("/system-instructions/active", h.GetActive)
	r.With(authz.RequireScope("write")).Post("/system-instructions", h.Create)
	r.With(authz.RequireScope("write")).Post("/system-instructions/{id}/activate", h.Activate)
}

func (h *Handler) RegisterTeacherRoutes(r chi.Router) {
	// Teachers have no access to system instructions
}

func (h *Handler) RegisterLearnerRoutes(r chi.Router) {
	// Learners have no access to system instructions
}

// ListAll lists all system instructions.
// @Summary List all system instructions
// @Description Get all system instructions
// @Tags System Instructions
// @Security OAuth2[read]
// @Success 200 {array} SystemInstruction "List of system instructions"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /system-instructions [get]
func (h *Handler) ListAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	instructions, err := h.service.ListAll(ctx)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, instructions)
}

// GetByID retrieves a system instruction by ID.
// @Summary Get system instruction by ID
// @Description Retrieve a specific system instruction by its UUID
// @Tags System Instructions
// @Security OAuth2[read]
// @Param id path string true "Instruction ID (UUID)"
// @Success 200 {object} SystemInstruction "System instruction details"
// @Failure 400 {object} map[string]string "Bad request - invalid ID format"
// @Failure 404 {object} map[string]string "Instruction not found"
// @Router /system-instructions/{id} [get]
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid instruction ID")
		return
	}

	ctx := r.Context()
	instruction, err := h.service.GetByID(ctx, id)
	if err != nil {
		render.Error(w, http.StatusNotFound, "Instruction not found")
		return
	}

	render.JSON(w, http.StatusOK, instruction)
}

// GetActive retrieves the active system instruction.
// @Summary Get active system instruction
// @Description Retrieve the currently active system instruction
// @Tags System Instructions
// @Security OAuth2[read]
// @Success 200 {object} SystemInstruction "Active system instruction"
// @Failure 404 {object} map[string]string "Active instruction not found"
// @Router /system-instructions/active [get]
func (h *Handler) GetActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	instruction, err := h.service.GetActive(ctx)
	if err != nil {
		render.Error(w, http.StatusNotFound, "Active instruction not found")
		return
	}

	render.JSON(w, http.StatusOK, instruction)
}

// Create creates a new system instruction.
// @Summary Create new system instruction
// @Description Create a new system instruction (immutable - cannot edit existing)
// @Tags System Instructions
// @Security OAuth2[write]
// @Accept json
// @Param request body CreateSystemInstructionRequest true "System instruction request"
// @Success 201 {object} SystemInstruction "Created system instruction"
// @Failure 400 {object} map[string]string "Bad request - invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /system-instructions [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateSystemInstructionRequest
	if err := render.DecodeJSON(r, &req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := r.Context()
	instruction, err := h.service.Create(ctx, req)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusCreated, instruction)
}

// Activate marks a system instruction as active.
// @Summary Activate system instruction
// @Description Mark a system instruction as active (deactivates other versions)
// @Tags System Instructions
// @Security OAuth2[write]
// @Param id path string true "Instruction ID (UUID)"
// @Success 200 {object} map[string]string "Activation status"
// @Failure 400 {object} map[string]string "Bad request - invalid ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /system-instructions/{id}/activate [post]
func (h *Handler) Activate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid instruction ID")
		return
	}

	ctx := r.Context()
	if err := h.service.Activate(ctx, id); err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"status": "activated"})
}
