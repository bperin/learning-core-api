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

// Create creates a schema template with an explicit version.
func (r *RepositoryImpl) Create(ctx context.Context, req CreateSchemaTemplateRequest) (*SchemaTemplate, error) {
	storeTemplate, err := r.queries.CreateSchemaTemplate(ctx, store.CreateSchemaTemplateParams{
		SchemaType:   req.SchemaType,
		Version:      req.Version,
		SchemaJson:   req.SchemaJSON,
		SubjectID:    utils.PtrToNullUUID(req.SubjectID),
		CurriculumID: utils.PtrToNullUUID(req.CurriculumID),
		IsActive:     req.IsActive,
		CreatedBy:    req.CreatedBy,
		LockedAt:     utils.SqlNullTime(req.LockedAt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create schema template: %w", err)
	}

	return toDomainSchemaTemplate(&storeTemplate), nil
}

// CreateVersion creates a new schema template version.
func (r *RepositoryImpl) CreateVersion(ctx context.Context, req CreateSchemaTemplateVersionRequest) (*SchemaTemplate, error) {
	storeTemplate, err := r.queries.CreateSchemaTemplateVersion(ctx, store.CreateSchemaTemplateVersionParams{
		SchemaType:   req.SchemaType,
		SchemaJson:   req.SchemaJSON,
		SubjectID:    utils.PtrToNullUUID(req.SubjectID),
		CurriculumID: utils.PtrToNullUUID(req.CurriculumID),
		IsActive:     req.IsActive,
		CreatedBy:    req.CreatedBy,
		LockedAt:     utils.SqlNullTime(req.LockedAt),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create schema template version: %w", err)
	}

	return toDomainSchemaTemplate(&storeTemplate), nil
}

// GetByID retrieves a schema template by ID.
func (r *RepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error) {
	storeTemplate, err := r.queries.GetSchemaTemplate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get schema template: %w", err)
	}

	return toDomainSchemaTemplate(&storeTemplate), nil
}

// GetByTypeAndVersion retrieves a schema template by type and version.
func (r *RepositoryImpl) GetByTypeAndVersion(ctx context.Context, schemaType string, version int32) (*SchemaTemplate, error) {
	storeTemplate, err := r.queries.GetSchemaTemplateByTypeAndVersion(ctx, store.GetSchemaTemplateByTypeAndVersionParams{
		SchemaType: schemaType,
		Version:    version,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get schema template by type/version: %w", err)
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

// ListBySubject lists schema templates by subject.
func (r *RepositoryImpl) ListBySubject(ctx context.Context, subjectID uuid.UUID) ([]*SchemaTemplate, error) {
	storeTemplates, err := r.queries.ListSchemaTemplatesBySubject(ctx, utils.UUIDToNullUUID(subjectID))
	if err != nil {
		return nil, fmt.Errorf("failed to list schema templates by subject: %w", err)
	}

	return toDomainSchemaTemplates(storeTemplates), nil
}

// ListByCurriculum lists schema templates by curriculum.
func (r *RepositoryImpl) ListByCurriculum(ctx context.Context, curriculumID uuid.UUID) ([]*SchemaTemplate, error) {
	storeTemplates, err := r.queries.ListSchemaTemplatesByCurriculum(ctx, utils.UUIDToNullUUID(curriculumID))
	if err != nil {
		return nil, fmt.Errorf("failed to list schema templates by curriculum: %w", err)
	}

	return toDomainSchemaTemplates(storeTemplates), nil
}

// Activate sets a schema template as active.
func (r *RepositoryImpl) Activate(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error) {
	storeTemplate, err := r.queries.ActivateSchemaTemplate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to activate schema template: %w", err)
	}

	return toDomainSchemaTemplate(&storeTemplate), nil
}

// Deactivate sets a schema template as inactive and locks it.
func (r *RepositoryImpl) Deactivate(ctx context.Context, id uuid.UUID) (*SchemaTemplate, error) {
	storeTemplate, err := r.queries.DeactivateSchemaTemplate(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to deactivate schema template: %w", err)
	}

	return toDomainSchemaTemplate(&storeTemplate), nil
}

// DeactivateOtherVersions deactivates other versions for a schema type.
func (r *RepositoryImpl) DeactivateOtherVersions(ctx context.Context, schemaType string, id uuid.UUID) error {
	if err := r.queries.DeactivateOtherSchemaTemplateVersions(ctx, store.DeactivateOtherSchemaTemplateVersionsParams{
		SchemaType: schemaType,
		ID:         id,
	}); err != nil {
		return fmt.Errorf("failed to deactivate other schema template versions: %w", err)
	}

	return nil
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
