package tenants

import (
	"context"
	"learning-core-api/internal/store"

	"github.com/google/uuid"
)

// Repository defines the interface for tenant operations
type Repository interface {
	Create(ctx context.Context, tenant Tenant) (*Tenant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
}

// repository implements the Repository interface
type repository struct {
	queries *store.Queries
}

// NewRepository creates a new tenant repository
func NewRepository(queries *store.Queries) Repository {
	return &repository{
		queries: queries,
	}
}

// Create creates a new tenant
func (r *repository) Create(ctx context.Context, tenant Tenant) (*Tenant, error) {
	params := store.CreateTenantParams{
		Name:     tenant.Name,
		IsActive: true,
	}

	dbTenant, err := r.queries.CreateTenant(ctx, params)
	if err != nil {
		return nil, err
	}

	tenant := toDomainTenant(dbTenant)
	return &tenant, nil
}

// GetByID retrieves a tenant by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	dbTenant, err := r.queries.GetTenantById(ctx, id)
	if err != nil {
		return nil, err
	}

	tenant := toDomainTenant(dbTenant)
	return &tenant, nil
}

func toDomainTenant(dbTenant store.Tenant) Tenant {
	return Tenant{
		ID:        dbTenant.ID,
		Name:      dbTenant.Name,
		IsActive:  dbTenant.IsActive,
		CreatedAt: dbTenant.CreatedAt,
	}
}
