package model_configs

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/utils"
)

// RepositoryImpl implements Repository using SQLC queries.
type RepositoryImpl struct {
	queries *store.Queries
}

// NewRepository creates a new model configs repository.
func NewRepository(queries *store.Queries) Repository {
	return &RepositoryImpl{queries: queries}
}

// Create creates a model config with an explicit version.
func (r *RepositoryImpl) Create(ctx context.Context, req CreateModelConfigRequest) (*ModelConfig, error) {
	isActive := false
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	storeConfig, err := r.queries.CreateModelConfig(ctx, store.CreateModelConfigParams{
		ModelName:   req.ModelName,
		Temperature: utils.SqlNullFloat64(req.Temperature),
		MaxTokens:   utils.SqlNullInt32(req.MaxTokens),
		TopP:        utils.SqlNullFloat64(req.TopP),
		TopK:        utils.SqlNullInt32(req.TopK),
		MimeType:    utils.SqlNullString(req.MimeType),
		IsActive:    isActive,
		CreatedBy:   req.CreatedBy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create model config: %w", err)
	}

	return toDomainModelConfigRow(&storeConfig), nil
}

// GetByID retrieves a model config by ID.
func (r *RepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*ModelConfig, error) {
	storeConfig, err := r.queries.GetModelConfig(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get model config: %w", err)
	}

	return toDomainModelConfig(&storeConfig), nil
}

// GetActive retrieves the active model config.
func (r *RepositoryImpl) GetActive(ctx context.Context) (*ModelConfig, error) {
	storeConfig, err := r.queries.GetActiveModelConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active model config: %w", err)
	}

	return toDomainModelConfig(&storeConfig), nil
}

// ListAll lists all model configs.
func (r *RepositoryImpl) ListAll(ctx context.Context) ([]*ModelConfig, error) {
	storeConfigs, err := r.queries.ListModelConfigs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list model configs: %w", err)
	}

	configs := make([]*ModelConfig, len(storeConfigs))
	for i, storeConfig := range storeConfigs {
		configs[i] = toDomainModelConfig(&storeConfig)
	}
	return configs, nil
}

// Activate marks a model config as active and deactivates others.
func (r *RepositoryImpl) Activate(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.ActivateModelConfig(ctx, id); err != nil {
		return fmt.Errorf("failed to activate model config: %w", err)
	}
	return nil
}

func toDomainModelConfig(storeConfig *store.ModelConfig) *ModelConfig {
	return &ModelConfig{
		ID:          storeConfig.ID,
		Version:     storeConfig.Version,
		ModelName:   storeConfig.ModelName,
		Temperature: utils.NullFloat64ToPtr(storeConfig.Temperature),
		MaxTokens:   utils.NullInt32ToPtr(storeConfig.MaxTokens),
		TopP:        utils.NullFloat64ToPtr(storeConfig.TopP),
		TopK:        utils.NullInt32ToPtr(storeConfig.TopK),
		MimeType:    utils.NullStringToPtr(storeConfig.MimeType),
		IsActive:    storeConfig.IsActive,
		CreatedBy:   storeConfig.CreatedBy,
		CreatedAt:   storeConfig.CreatedAt,
	}
}

func toDomainModelConfigRow(storeConfig *store.CreateModelConfigRow) *ModelConfig {
	return &ModelConfig{
		ID:          storeConfig.ID,
		Version:     storeConfig.Version,
		ModelName:   storeConfig.ModelName,
		Temperature: utils.NullFloat64ToPtr(storeConfig.Temperature),
		MaxTokens:   utils.NullInt32ToPtr(storeConfig.MaxTokens),
		TopP:        utils.NullFloat64ToPtr(storeConfig.TopP),
		TopK:        utils.NullInt32ToPtr(storeConfig.TopK),
		MimeType:    utils.NullStringToPtr(storeConfig.MimeType),
		IsActive:    storeConfig.IsActive,
		CreatedBy:   storeConfig.CreatedBy,
		CreatedAt:   storeConfig.CreatedAt,
	}
}
