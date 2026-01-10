package chunking_configs

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

// NewRepository creates a new chunking configs repository.
func NewRepository(queries *store.Queries) Repository {
	return &RepositoryImpl{queries: queries}
}

// Create creates a chunking config.
func (r *RepositoryImpl) Create(ctx context.Context, req CreateChunkingConfigRequest) (*ChunkingConfig, error) {
	isActive := false
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	storeConfig, err := r.queries.CreateChunkingConfig(ctx, store.CreateChunkingConfigParams{
		ChunkSize:    req.ChunkSize,
		ChunkOverlap: req.ChunkOverlap,
		IsActive:     isActive,
		CreatedBy:    req.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create chunking config: %w", err)
	}

	if isActive {
		if err := r.Activate(ctx, storeConfig.ID); err != nil {
			return nil, err
		}
		return r.GetByID(ctx, storeConfig.ID)
	}

	return toDomainChunkingConfig(&storeConfig), nil
}

// GetByID retrieves a chunking config by ID.
func (r *RepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*ChunkingConfig, error) {
	storeConfig, err := r.queries.GetChunkingConfig(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get chunking config: %w", err)
	}

	return toDomainChunkingConfig(&storeConfig), nil
}

// GetActive retrieves the active chunking config.
func (r *RepositoryImpl) GetActive(ctx context.Context) (*ChunkingConfig, error) {
	storeConfig, err := r.queries.GetActiveChunkingConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active chunking config: %w", err)
	}

	return toDomainChunkingConfig(&storeConfig), nil
}

// ListAll lists all chunking configs.
func (r *RepositoryImpl) ListAll(ctx context.Context) ([]*ChunkingConfig, error) {
	storeConfigs, err := r.queries.ListChunkingConfigs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list chunking configs: %w", err)
	}

	configs := make([]*ChunkingConfig, len(storeConfigs))
	for i, storeConfig := range storeConfigs {
		configs[i] = toDomainChunkingConfig(&storeConfig)
	}
	return configs, nil
}

// Activate marks a chunking config as active and deactivates others.
func (r *RepositoryImpl) Activate(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.ActivateChunkingConfig(ctx, id); err != nil {
		return fmt.Errorf("failed to activate chunking config: %w", err)
	}
	if err := r.queries.DeactivateOtherChunkingConfigs(ctx, id); err != nil {
		return fmt.Errorf("failed to deactivate other chunking configs: %w", err)
	}
	return nil
}

func toDomainChunkingConfig(storeConfig *store.ChunkingConfig) *ChunkingConfig {
	return &ChunkingConfig{
		ID:           storeConfig.ID,
		Version:      storeConfig.Version,
		ChunkSize:    storeConfig.ChunkSize,
		ChunkOverlap: storeConfig.ChunkOverlap,
		IsActive:     storeConfig.IsActive,
		CreatedBy:    storeConfig.CreatedBy,
		CreatedAt:    storeConfig.CreatedAt,
	}
}
