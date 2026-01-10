package prompt_templates

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"

	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/utils"
)

// RepositoryImpl implements Repository using SQLC queries.
type RepositoryImpl struct {
	queries *store.Queries
}

// NewRepository creates a new prompt templates repository.
func NewRepository(queries *store.Queries) Repository {
	return &RepositoryImpl{queries: queries}
}

// Create creates a prompt template with an explicit version.
func (r *RepositoryImpl) Create(ctx context.Context, req CreatePromptTemplateRequest) (*PromptTemplate, error) {
	storeTemplate, err := r.queries.CreatePromptTemplate(ctx, store.CreatePromptTemplateParams{
		Key:         req.Key,
		Version:     req.Version,
		IsActive:    req.IsActive,
		Title:       req.Title,
		Description: utils.SqlNullString(req.Description),
		Template:    req.Template,
		Metadata:    utils.ToNullRawMessage(req.Metadata),
		CreatedBy:   utils.SqlNullString(req.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt template: %w", err)
	}

	return toDomainPromptTemplateCreate(&storeTemplate), nil
}

// CreateVersion creates a new prompt template version.
func (r *RepositoryImpl) CreateVersion(ctx context.Context, req CreatePromptTemplateVersionRequest) (*PromptTemplate, error) {
	storeTemplate, err := r.queries.CreateNewVersion(ctx, store.CreateNewVersionParams{
		Key:         req.Key,
		IsActive:    req.IsActive,
		Title:       req.Title,
		Description: utils.SqlNullString(req.Description),
		Template:    req.Template,
		Metadata:    utils.ToNullRawMessage(req.Metadata),
		CreatedBy:   utils.SqlNullString(req.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt template version: %w", err)
	}

	return toDomainPromptTemplateVersion(&storeTemplate), nil
}

// GetByID retrieves a prompt template by ID.
func (r *RepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*PromptTemplate, error) {
	storeTemplate, err := r.queries.GetPromptTemplate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt template: %w", err)
	}

	return toDomainPromptTemplate(&storeTemplate), nil
}

// GetActiveByKey retrieves the active prompt template for a key.
func (r *RepositoryImpl) GetActiveByKey(ctx context.Context, key string) (*PromptTemplate, error) {
	storeTemplate, err := r.queries.GetPromptTemplateByKey(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get active prompt template: %w", err)
	}

	return toDomainPromptTemplate(&storeTemplate), nil
}

// GetByKeyAndVersion retrieves a prompt template by key and version.
func (r *RepositoryImpl) GetByKeyAndVersion(ctx context.Context, key string, version int32) (*PromptTemplate, error) {
	storeTemplate, err := r.queries.GetPromptTemplateByKeyAndVersion(ctx, store.GetPromptTemplateByKeyAndVersionParams{
		Key:     key,
		Version: version,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get prompt template by key/version: %w", err)
	}

	return toDomainPromptTemplate(&storeTemplate), nil
}

// Activate marks a prompt template as active.
func (r *RepositoryImpl) Activate(ctx context.Context, id uuid.UUID) (*PromptTemplate, error) {
	storeTemplate, err := r.queries.ActivatePromptTemplate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to activate prompt template: %w", err)
	}

	return toDomainPromptTemplateActivate(&storeTemplate), nil
}

// Deactivate marks a prompt template as inactive.
func (r *RepositoryImpl) Deactivate(ctx context.Context, id uuid.UUID) (*PromptTemplate, error) {
	storeTemplate, err := r.queries.DeactivatePromptTemplate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to deactivate prompt template: %w", err)
	}

	return toDomainPromptTemplate(&storeTemplate), nil
}

// DeactivateOtherVersions deactivates other versions for a prompt template key.
func (r *RepositoryImpl) DeactivateOtherVersions(ctx context.Context, key string, id uuid.UUID) error {
	if err := r.queries.DeactivateOtherVersions(ctx, store.DeactivateOtherVersionsParams{
		Key: key,
		ID:  id,
	}); err != nil {
		return fmt.Errorf("failed to deactivate other prompt template versions: %w", err)
	}
	return nil
}

func toDomainPromptTemplate(storeTemplate *store.PromptTemplate) *PromptTemplate {
	var metadata json.RawMessage
	if storeTemplate.Metadata.Valid {
		metadata = storeTemplate.Metadata.RawMessage
	}

	return &PromptTemplate{
		ID:          storeTemplate.ID,
		Key:         storeTemplate.Key,
		Version:     storeTemplate.Version,
		IsActive:    storeTemplate.IsActive,
		Title:       storeTemplate.Title,
		Description: utils.NullStringToPtr(storeTemplate.Description),
		Template:    storeTemplate.Template,
		Metadata:    metadata,
		CreatedBy:   utils.NullStringToPtr(storeTemplate.CreatedBy),
		CreatedAt:   storeTemplate.CreatedAt,
		UpdatedAt:   storeTemplate.UpdatedAt,
	}
}

func toDomainPromptTemplateCreate(storeTemplate *store.CreatePromptTemplateRow) *PromptTemplate {
	return toDomainPromptTemplateFields(
		storeTemplate.ID,
		storeTemplate.Key,
		storeTemplate.Version,
		storeTemplate.IsActive,
		storeTemplate.Title,
		storeTemplate.Description,
		storeTemplate.Template,
		storeTemplate.Metadata,
		storeTemplate.CreatedBy,
		storeTemplate.CreatedAt,
		storeTemplate.UpdatedAt,
	)
}

func toDomainPromptTemplateVersion(storeTemplate *store.CreateNewVersionRow) *PromptTemplate {
	return toDomainPromptTemplateFields(
		storeTemplate.ID,
		storeTemplate.Key,
		storeTemplate.Version,
		storeTemplate.IsActive,
		storeTemplate.Title,
		storeTemplate.Description,
		storeTemplate.Template,
		storeTemplate.Metadata,
		storeTemplate.CreatedBy,
		storeTemplate.CreatedAt,
		storeTemplate.UpdatedAt,
	)
}

func toDomainPromptTemplateActivate(storeTemplate *store.ActivatePromptTemplateRow) *PromptTemplate {
	return toDomainPromptTemplateFields(
		storeTemplate.ID,
		storeTemplate.Key,
		storeTemplate.Version,
		storeTemplate.IsActive,
		storeTemplate.Title,
		storeTemplate.Description,
		storeTemplate.Template,
		storeTemplate.Metadata,
		storeTemplate.CreatedBy,
		storeTemplate.CreatedAt,
		storeTemplate.UpdatedAt,
	)
}

func toDomainPromptTemplateFields(
	id uuid.UUID,
	key string,
	version int32,
	isActive bool,
	title string,
	description sql.NullString,
	templateText string,
	metadataRaw pqtype.NullRawMessage,
	createdBy sql.NullString,
	createdAt time.Time,
	updatedAt time.Time,
) *PromptTemplate {
	var metadata json.RawMessage
	if metadataRaw.Valid {
		metadata = metadataRaw.RawMessage
	}

	return &PromptTemplate{
		ID:          id,
		Key:         key,
		Version:     version,
		IsActive:    isActive,
		Title:       title,
		Description: utils.NullStringToPtr(description),
		Template:    templateText,
		Metadata:    metadata,
		CreatedBy:   utils.NullStringToPtr(createdBy),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
