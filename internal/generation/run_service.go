package generation

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// CreateGenerationRun creates a new generation run with business logic validation.
func (s *service) CreateGenerationRun(ctx context.Context, req CreateGenerationRunRequest) (*GenerationRun, error) {
	if req.ModuleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	if req.AgentName == "" {
		return nil, errors.New("agent name is required")
	}

	if req.Model == "" {
		return nil, errors.New("model is required")
	}

	if req.StoreName == "" {
		return nil, errors.New("store name is required")
	}

	run := GenerationRun{
		ModuleID:       req.ModuleID,
		AgentName:      req.AgentName,
		AgentVersion:   req.AgentVersion,
		Model:          req.Model,
		ModelParams:    req.ModelParams,
		PromptID:       req.PromptID,
		StoreName:      req.StoreName,
		MetadataFilter: req.MetadataFilter,
		Status:         RunStatusPending,
		InputPayload:   req.InputPayload,
		OutputPayload:  nil,
		Error:          nil,
		StartedAt:      nil,
		FinishedAt:     nil,
	}

	return s.repo.CreateGenerationRun(ctx, run)
}

// GetGenerationRunByID retrieves a generation run by ID.
func (s *service) GetGenerationRunByID(ctx context.Context, id uuid.UUID) (*GenerationRun, error) {
	if id == uuid.Nil {
		return nil, errors.New("generation run ID is required")
	}

	return s.repo.GetGenerationRunByID(ctx, id)
}

// ListGenerationRunsByModule retrieves all generation runs for a module.
func (s *service) ListGenerationRunsByModule(ctx context.Context, moduleID uuid.UUID) ([]GenerationRun, error) {
	if moduleID == uuid.Nil {
		return nil, errors.New("module ID is required")
	}

	return s.repo.ListGenerationRunsByModule(ctx, moduleID)
}

// UpdateGenerationRun updates a generation run with business logic validation.
func (s *service) UpdateGenerationRun(ctx context.Context, id uuid.UUID, status *RunStatus, outputPayload json.RawMessage, error json.RawMessage, startedAt, finishedAt *time.Time) error {
	if id == uuid.Nil {
		return errors.New("generation run ID is required")
	}

	_, err := s.repo.GetGenerationRunByID(ctx, id)
	if err != nil {
		return err
	}

	if status != nil {
		currentRun, err := s.repo.GetGenerationRunByID(ctx, id)
		if err != nil {
			return err
		}

		if !isValidStatusTransition(currentRun.Status, *status) {
			return errors.New("invalid status transition")
		}
	}

	return s.repo.UpdateGenerationRun(ctx, id, status, outputPayload, error, startedAt, finishedAt)
}

// isValidStatusTransition checks if a status transition is valid.
func isValidStatusTransition(from, to RunStatus) bool {
	validTransitions := map[RunStatus][]RunStatus{
		RunStatusPending:   {RunStatusRunning, RunStatusFailed},
		RunStatusRunning:   {RunStatusCompleted, RunStatusFailed},
		RunStatusCompleted: {},
		RunStatusFailed:    {},
	}

	for _, validTo := range validTransitions[from] {
		if to == validTo {
			return true
		}
	}

	return false
}

// DeleteGenerationRun deletes a generation run.
func (s *service) DeleteGenerationRun(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("generation run ID is required")
	}

	_, err := s.repo.GetGenerationRunByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteGenerationRun(ctx, id)
}
