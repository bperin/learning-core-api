package schema_templates

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

// NewRepository creates a new schema templates repository.
func NewRepository(queries *store.Queries) Repository {
	return &RepositoryImpl{queries: queries}
}

// Create creates a schema template with an auto-incremented version.
func (r *RepositoryImpl) Create(ctx context.Context, req CreateSchemaTemplateRequest) (*SchemaTemplate, error) {
	isActive := false
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	storeTemplate, err := r.queries.CreateSchemaTemplate(ctx, store.CreateSchemaTemplateParams{
		SchemaType:   req.SchemaType,
		SchemaJson:   req.SchemaJSON,
		SubjectID:    utils.PtrToNullUUID(req.SubjectID),
		CurriculumID: utils.PtrToNullUUID(req.CurriculumID),
		IsActive:     isActive,
		CreatedBy:    req.CreatedBy,
		LockedAt:     utils.SqlNullTime(req.LockedAt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create schema template: %w", err)
	}

	return toDomainSchemaTemplateCreate(&storeTemplate), nil
}

// GetByID retrieves a schema template by ID.
func (r *RepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error) {
	storeTemplate, err := r.queries.GetSchemaTemplate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema template: %w", err)
	}

	return toDomainSchemaTemplate(&storeTemplate), nil
}

// GetActiveByType retrieves the active schema template for a type.
func (r *RepositoryImpl) GetActiveByType(ctx context.Context, schemaType string) (*SchemaTemplate, error) {
	storeTemplate, err := r.queries.GetActiveSchemaTemplateByType(ctx, schemaType)
	if err != nil {
		return nil, fmt.Errorf("failed to get active schema template: %w", err)
	}

	return toDomainSchemaTemplate(&storeTemplate), nil
}

// ListByType lists schema templates by type.
func (r *RepositoryImpl) ListByType(ctx context.Context, schemaType string) ([]*SchemaTemplate, error) {
	storeTemplates, err := r.queries.ListSchemaTemplatesByType(ctx, schemaType)
	if err != nil {
		return nil, fmt.Errorf("failed to list schema templates by type: %w", err)
	}

	return toDomainSchemaTemplates(storeTemplates), nil
}

// ListActive lists active schema templates.
func (r *RepositoryImpl) ListActive(ctx context.Context) ([]*SchemaTemplate, error) {
	storeTemplates, err := r.queries.ListActiveSchemaTemplates(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list active schema templates: %w", err)
	}

	return toDomainSchemaTemplates(storeTemplates), nil
}

// Activate sets a schema template as active, deactivating other versions.
func (r *RepositoryImpl) Activate(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error) {
	storeTemplate, err := r.queries.ActivateSchemaTemplate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to activate schema template: %w", err)
	}

	return toDomainSchemaTemplateRow(&storeTemplate), nil
}

func toDomainSchemaTemplates(storeTemplates []store.SchemaTemplate) []*SchemaTemplate {
	templates := make([]*SchemaTemplate, len(storeTemplates))
	for i, storeTemplate := range storeTemplates {
		templates[i] = toDomainSchemaTemplate(&storeTemplate)
	}
	return templates
}

func toDomainSchemaTemplate(storeTemplate *store.SchemaTemplate) *SchemaTemplate {
	return &SchemaTemplate{
		ID:           storeTemplate.ID,
		SchemaType:   storeTemplate.SchemaType,
		Version:      storeTemplate.Version,
		SchemaJSON:   storeTemplate.SchemaJson,
		SubjectID:    utils.NullUUIDToPtr(storeTemplate.SubjectID),
		CurriculumID: utils.NullUUIDToPtr(storeTemplate.CurriculumID),
		IsActive:     storeTemplate.IsActive,
		CreatedBy:    storeTemplate.CreatedBy,
		CreatedAt:    storeTemplate.CreatedAt,
		LockedAt:     utils.NullTimeToPtr(storeTemplate.LockedAt),
	}
}

func toDomainSchemaTemplateRow(storeTemplate *store.ActivateSchemaTemplateRow) *SchemaTemplate {
	return &SchemaTemplate{
		ID:           storeTemplate.ID,
		SchemaType:   storeTemplate.SchemaType,
		Version:      storeTemplate.Version,
		SchemaJSON:   storeTemplate.SchemaJson,
		SubjectID:    utils.NullUUIDToPtr(storeTemplate.SubjectID),
		CurriculumID: utils.NullUUIDToPtr(storeTemplate.CurriculumID),
		IsActive:     storeTemplate.IsActive,
		CreatedBy:    storeTemplate.CreatedBy,
		CreatedAt:    storeTemplate.CreatedAt,
		LockedAt:     utils.NullTimeToPtr(storeTemplate.LockedAt),
	}
}

func toDomainSchemaTemplateCreate(storeTemplate *store.CreateSchemaTemplateRow) *SchemaTemplate {
	return &SchemaTemplate{
		ID:           storeTemplate.ID,
		SchemaType:   storeTemplate.SchemaType,
		Version:      storeTemplate.Version,
		SchemaJSON:   storeTemplate.SchemaJson,
		SubjectID:    utils.NullUUIDToPtr(storeTemplate.SubjectID),
		CurriculumID: utils.NullUUIDToPtr(storeTemplate.CurriculumID),
		IsActive:     storeTemplate.IsActive,
		CreatedBy:    storeTemplate.CreatedBy,
		CreatedAt:    storeTemplate.CreatedAt,
		LockedAt:     utils.NullTimeToPtr(storeTemplate.LockedAt),
	}
}
