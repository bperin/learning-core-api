package filesearch

import (
	"context"
	"database/sql"
	"encoding/json"
	"learning-core-api/internal/store"

	"github.com/google/uuid"
)

// Repository defines storage operations for file search stores.
type Repository interface {
	Create(ctx context.Context, store Store) (*Store, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Store, error)
	GetBySubjectID(ctx context.Context, subjectID uuid.UUID) (*Store, error)
	Update(ctx context.Context, id uuid.UUID, displayName string, chunkingConfig json.RawMessage) (*Store, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	queries *store.Queries
}

// NewRepository creates a new file search store repository.
func NewRepository(queries *store.Queries) Repository {
	return &repository{queries: queries}
}

func (r *repository) Create(ctx context.Context, fsStore Store) (*Store, error) {
	displayName := sql.NullString{Valid: false}
	if fsStore.DisplayName != "" {
		displayName = sql.NullString{String: fsStore.DisplayName, Valid: true}
	}

	dbStore, err := r.queries.CreateFileSearchStore(ctx, store.CreateFileSearchStoreParams{
		SubjectID:      fsStore.SubjectID,
		StoreName:      fsStore.StoreName,
		DisplayName:    displayName,
		ChunkingConfig: fsStore.ChunkingConfig,
	})
	if err != nil {
		return nil, err
	}

	store := toDomainStore(dbStore)
	return &store, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Store, error) {
	dbStore, err := r.queries.GetFileSearchStore(ctx, id)
	if err != nil {
		return nil, err
	}

	store := toDomainStore(dbStore)
	return &store, nil
}

func (r *repository) GetBySubjectID(ctx context.Context, subjectID uuid.UUID) (*Store, error) {
	dbStore, err := r.queries.GetFileSearchStoreBySubjectID(ctx, subjectID)
	if err != nil {
		return nil, err
	}

	store := toDomainStore(dbStore)
	return &store, nil
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, displayName string, chunkingConfig json.RawMessage) (*Store, error) {
	displayNameValue := sql.NullString{Valid: false}
	if displayName != "" {
		displayNameValue = sql.NullString{String: displayName, Valid: true}
	}

	dbStore, err := r.queries.UpdateFileSearchStore(ctx, store.UpdateFileSearchStoreParams{
		ID:             id,
		DisplayName:    displayNameValue,
		ChunkingConfig: chunkingConfig,
	})
	if err != nil {
		return nil, err
	}

	store := toDomainStore(dbStore)
	return &store, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteFileSearchStore(ctx, id)
}

func toDomainStore(dbStore store.FileSearchStore) Store {
	displayName := dbStore.DisplayName.String
	if !dbStore.DisplayName.Valid {
		displayName = dbStore.StoreName
	}

	return Store{
		ID:             dbStore.ID,
		SubjectID:      dbStore.SubjectID,
		StoreName:      dbStore.StoreName,
		DisplayName:    displayName,
		ChunkingConfig: dbStore.ChunkingConfig,
		CreatedAt:      dbStore.CreatedAt,
	}
}
