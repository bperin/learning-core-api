package filesearch

import (
	"context"
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
	dbStore, err := r.queries.CreateFileSearchStore(ctx, store.CreateFileSearchStoreParams{
		SubjectID:      fsStore.SubjectID,
		StoreName:      fsStore.StoreName,
		DisplayName:    fsStore.DisplayName,
		ChunkingConfig: fsStore.ChunkingConfig,
	})
	if err != nil {
		return nil, err
	}

	return &Store{
		ID:             dbStore.ID,
		SubjectID:      dbStore.SubjectID,
		StoreName:      dbStore.StoreName,
		DisplayName:    dbStore.DisplayName,
		ChunkingConfig: dbStore.ChunkingConfig,
		CreatedAt:      dbStore.CreatedAt,
	}, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Store, error) {
	dbStore, err := r.queries.GetFileSearchStore(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Store{
		ID:             dbStore.ID,
		SubjectID:      dbStore.SubjectID,
		StoreName:      dbStore.StoreName,
		DisplayName:    dbStore.DisplayName,
		ChunkingConfig: dbStore.ChunkingConfig,
		CreatedAt:      dbStore.CreatedAt,
	}, nil
}

func (r *repository) GetBySubjectID(ctx context.Context, subjectID uuid.UUID) (*Store, error) {
	dbStore, err := r.queries.GetFileSearchStoreBySubjectID(ctx, subjectID)
	if err != nil {
		return nil, err
	}

	return &Store{
		ID:             dbStore.ID,
		SubjectID:      dbStore.SubjectID,
		StoreName:      dbStore.StoreName,
		DisplayName:    dbStore.DisplayName,
		ChunkingConfig: dbStore.ChunkingConfig,
		CreatedAt:      dbStore.CreatedAt,
	}, nil
}

func (r *repository) Update(ctx context.Context, id uuid.UUID, displayName string, chunkingConfig json.RawMessage) (*Store, error) {
	dbStore, err := r.queries.UpdateFileSearchStore(ctx, store.UpdateFileSearchStoreParams{
		ID:             id,
		DisplayName:    displayName,
		ChunkingConfig: chunkingConfig,
	})
	if err != nil {
		return nil, err
	}

	return &Store{
		ID:             dbStore.ID,
		SubjectID:      dbStore.SubjectID,
		StoreName:      dbStore.StoreName,
		DisplayName:    dbStore.DisplayName,
		ChunkingConfig: dbStore.ChunkingConfig,
		CreatedAt:      dbStore.CreatedAt,
	}, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteFileSearchStore(ctx, id)
}
