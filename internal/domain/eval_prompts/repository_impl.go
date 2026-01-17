package eval_prompts

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"learning-core-api/internal/persistance/store"
)

// RepositoryImpl implements the Repository interface using SQLC
type RepositoryImpl struct {
	queries *store.Queries
}

// NewRepository creates a new eval prompt repository
func NewRepository(queries *store.Queries) Repository {
	return &RepositoryImpl{
		queries: queries,
	}
}

// GetByID retrieves an eval prompt by ID
func (r *RepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*EvalPrompt, error) {
	prompt, err := r.queries.GetEvalPrompt(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get eval prompt: %w", err)
	}

	return r.mapToEvalPrompt(&prompt), nil
}

// GetActiveByType retrieves the active eval prompt for a given eval type
func (r *RepositoryImpl) GetActiveByType(ctx context.Context, evalType string) (*EvalPrompt, error) {
	prompt, err := r.queries.GetActiveEvalPrompt(ctx, evalType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get active eval prompt: %w", err)
	}

	return r.mapToEvalPrompt(&prompt), nil
}

// GetByTypeAndVersion retrieves a specific version of an eval prompt
func (r *RepositoryImpl) GetByTypeAndVersion(ctx context.Context, evalType string, version int32) (*EvalPrompt, error) {
	prompt, err := r.queries.GetEvalPromptByVersion(ctx, store.GetEvalPromptByVersionParams{
		EvalType: evalType,
		Version:  version,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get eval prompt by version: %w", err)
	}

	return r.mapToEvalPrompt(&prompt), nil
}

// ListByType retrieves all versions of an eval prompt type
func (r *RepositoryImpl) ListByType(ctx context.Context, evalType string, limit int32, offset int32) ([]*EvalPrompt, error) {
	prompts, err := r.queries.ListEvalPrompts(ctx, store.ListEvalPromptsParams{
		EvalType: evalType,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list eval prompts: %w", err)
	}

	result := make([]*EvalPrompt, len(prompts))
	for i, p := range prompts {
		result[i] = r.mapToEvalPrompt(&p)
	}

	return result, nil
}

// Create creates a new eval prompt
func (r *RepositoryImpl) Create(ctx context.Context, req *CreateEvalPromptRequest) (*EvalPrompt, error) {
	// Get the next version number
	latestVersion, err := r.GetLatestVersion(ctx, req.EvalType)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest version: %w", err)
	}

	nextVersion := latestVersion + 1

	prompt, err := r.queries.CreateEvalPrompt(ctx, store.CreateEvalPromptParams{
		EvalType:    req.EvalType,
		Version:     nextVersion,
		PromptText:  req.PromptText,
		Description: toNullString(req.Description),
		IsActive:    sql.NullBool{Bool: true, Valid: true},
		CreatedBy:   toNullUUID(req.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create eval prompt: %w", err)
	}

	return r.mapToEvalPrompt(prompt), nil
}

// Activate activates a specific eval prompt version
func (r *RepositoryImpl) Activate(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.ActivateEvalPrompt(ctx, id); err != nil {
		return fmt.Errorf("failed to activate eval prompt: %w", err)
	}

	return nil
}

// Deactivate deactivates a specific eval prompt version
func (r *RepositoryImpl) Deactivate(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.DeactivateEvalPrompt(ctx, id); err != nil {
		return fmt.Errorf("failed to deactivate eval prompt: %w", err)
	}

	return nil
}

// GetLatestVersion gets the latest version number for an eval type
func (r *RepositoryImpl) GetLatestVersion(ctx context.Context, evalType string) (int32, error) {
	// This query returns coalesce(max(version), 0) so it should never be nil
	// For now, we'll handle it gracefully
	result, err := r.queries.GetLatestEvalPromptVersion(ctx, evalType)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest version: %w", err)
	}

	// The result should be int32 or similar
	if v, ok := result.(int32); ok {
		return v, nil
	}

	return 0, nil
}

// mapToEvalPrompt converts a store.EvalPrompt to a domain EvalPrompt
func (r *RepositoryImpl) mapToEvalPrompt(prompt *store.EvalPrompt) *EvalPrompt {
	if prompt == nil {
		return nil
	}

	var description *string
	if prompt.Description.Valid {
		description = &prompt.Description.String
	}

	var createdBy *uuid.UUID
	if prompt.CreatedBy.Valid {
		createdBy = &prompt.CreatedBy.UUID
	}

	var createdAt, updatedAt sql.NullTime
	if prompt.CreatedAt.Valid {
		createdAt = prompt.CreatedAt
	}
	if prompt.UpdatedAt.Valid {
		updatedAt = prompt.UpdatedAt
	}

	return &EvalPrompt{
		ID:          prompt.ID,
		EvalType:    prompt.EvalType,
		Version:     prompt.Version,
		PromptText:  prompt.PromptText,
		Description: description,
		IsActive:    prompt.IsActive.Bool,
		CreatedBy:   createdBy,
		CreatedAt:   createdAt.Time,
		UpdatedAt:   updatedAt.Time,
	}
}

// toNullString converts a string pointer to sql.NullString
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// toNullUUID converts a uuid.UUID pointer to uuid.NullUUID
func toNullUUID(u *uuid.UUID) uuid.NullUUID {
	if u == nil {
		return uuid.NullUUID{Valid: false}
	}
	return uuid.NullUUID{UUID: *u, Valid: true}
}
