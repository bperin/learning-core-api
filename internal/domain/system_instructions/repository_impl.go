package system_instructions

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"learning-core-api/internal/persistance/store"
)

// RepositoryImpl implements Repository using SQLC queries.
type RepositoryImpl struct {
	queries *store.Queries
}

// NewRepository creates a new system instructions repository.
func NewRepository(queries *store.Queries) Repository {
	return &RepositoryImpl{queries: queries}
}

// Create creates a system instruction.
func (r *RepositoryImpl) Create(ctx context.Context, req CreateSystemInstructionRequest) (*SystemInstruction, error) {
	isActive := false
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	storeInstruction, err := r.queries.CreateSystemInstruction(ctx, store.CreateSystemInstructionParams{
		Text:      req.Text,
		IsActive:  isActive,
		CreatedBy: req.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create system instruction: %w", err)
	}

	if isActive {
		if err := r.Activate(ctx, storeInstruction.ID); err != nil {
			return nil, err
		}
		return r.GetByID(ctx, storeInstruction.ID)
	}

	return toDomainSystemInstruction(&storeInstruction), nil
}

// GetByID retrieves a system instruction by ID.
func (r *RepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*SystemInstruction, error) {
	storeInstruction, err := r.queries.GetSystemInstruction(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get system instruction: %w", err)
	}

	return toDomainSystemInstruction(&storeInstruction), nil
}

// GetActive retrieves the active system instruction.
func (r *RepositoryImpl) GetActive(ctx context.Context) (*SystemInstruction, error) {
	storeInstruction, err := r.queries.GetActiveSystemInstruction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active system instruction: %w", err)
	}

	return toDomainSystemInstruction(&storeInstruction), nil
}

// ListAll lists all system instructions.
func (r *RepositoryImpl) ListAll(ctx context.Context) ([]*SystemInstruction, error) {
	storeInstructions, err := r.queries.ListSystemInstructions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list system instructions: %w", err)
	}

	instructions := make([]*SystemInstruction, len(storeInstructions))
	for i, storeInstruction := range storeInstructions {
		instructions[i] = toDomainSystemInstruction(&storeInstruction)
	}
	return instructions, nil
}

// Activate marks a system instruction as active and deactivates others.
func (r *RepositoryImpl) Activate(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.ActivateSystemInstruction(ctx, id); err != nil {
		return fmt.Errorf("failed to activate system instruction: %w", err)
	}
	if err := r.queries.DeactivateOtherSystemInstructions(ctx, id); err != nil {
		return fmt.Errorf("failed to deactivate other system instructions: %w", err)
	}
	return nil
}

func toDomainSystemInstruction(storeInstruction *store.SystemInstruction) *SystemInstruction {
	return &SystemInstruction{
		ID:        storeInstruction.ID,
		Version:   storeInstruction.Version,
		Text:      storeInstruction.Text,
		IsActive:  storeInstruction.IsActive,
		CreatedBy: storeInstruction.CreatedBy,
		CreatedAt: storeInstruction.CreatedAt,
	}
}
