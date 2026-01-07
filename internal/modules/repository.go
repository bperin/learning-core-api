package modules

import (
	"context"
	"database/sql"
	"learning-core-api/internal/store"

	"github.com/google/uuid"
)

// Repository defines the interface for module operations
type Repository interface {
	Create(ctx context.Context, module Module) (*Module, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Module, error)
	GetByName(ctx context.Context, tenantID uuid.UUID, name string) (*Module, error)
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]Module, error)
	Update(ctx context.Context, id uuid.UUID, name, description *string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// repository implements the Repository interface
type repository struct {
	queries *store.Queries
}

// NewRepository creates a new module repository
func NewRepository(queries *store.Queries) Repository {
	return &repository{
		queries: queries,
	}
}

// Create creates a new module
func (r *repository) Create(ctx context.Context, module Module) (*Module, error) {
	params := store.CreateModuleParams{
		TenantID: module.TenantID,
		Title:    module.Title,
		Name:     sql.NullString{String: module.Name, Valid: true}, // Use Name field for Name column
	}

	if module.Description != "" {
		params.Description = sql.NullString{String: module.Description, Valid: true}
	} else {
		params.Description = sql.NullString{String: "", Valid: false}
	}

	dbModule, err := r.queries.CreateModule(ctx, params)
	if err != nil {
		return nil, err
	}

	description := ""
	if dbModule.Description.Valid {
		description = dbModule.Description.String
	}

	name := ""
	if dbModule.Name.Valid {
		name = dbModule.Name.String
	}

	return &Module{
		ID:          dbModule.ID,
		TenantID:    dbModule.TenantID,
		Title:       dbModule.Title,
		Name:        name, // Use Name field
		Description: description,
		CreatedAt:   dbModule.CreatedAt,
	}, nil
}

// GetByID retrieves a module by ID
func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Module, error) {
	dbModule, err := r.queries.GetModule(ctx, id)
	if err != nil {
		return nil, err
	}

	description := ""
	if dbModule.Description.Valid {
		description = dbModule.Description.String
	}

	name := ""
	if dbModule.Name.Valid {
		name = dbModule.Name.String
	}

	return &Module{
		ID:          dbModule.ID,
		TenantID:    dbModule.TenantID,
		Title:       dbModule.Title,
		Name:        name,
		Description: description,
		CreatedAt:   dbModule.CreatedAt,
	}, nil
}

// GetByName retrieves a module by name and tenant ID
func (r *repository) GetByName(ctx context.Context, tenantID uuid.UUID, name string) (*Module, error) {
	dbModule, err := r.queries.GetModuleByName(ctx, store.GetModuleByNameParams{
		TenantID: tenantID,
		Name:     sql.NullString{String: name, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	description := ""
	if dbModule.Description.Valid {
		description = dbModule.Description.String
	}

	moduleName := ""
	if dbModule.Name.Valid {
		moduleName = dbModule.Name.String
	}

	return &Module{
		ID:          dbModule.ID,
		TenantID:    dbModule.TenantID,
		Title:       dbModule.Title,
		Name:        moduleName,
		Description: description,
		CreatedAt:   dbModule.CreatedAt,
	}, nil
}

// ListByTenant retrieves all modules for a tenant
func (r *repository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]Module, error) {
	dbModules, err := r.queries.ListModulesByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	modules := make([]Module, len(dbModules))
	for i, dbModule := range dbModules {
		description := ""
		if dbModule.Description.Valid {
			description = dbModule.Description.String
		}

		name := ""
		if dbModule.Name.Valid {
			name = dbModule.Name.String
		}

		modules[i] = Module{
			ID:          dbModule.ID,
			TenantID:    dbModule.TenantID,
			Title:       dbModule.Title,
			Name:        name,
			Description: description,
			CreatedAt:   dbModule.CreatedAt,
		}
	}

	return modules, nil
}

// Update updates a module
func (r *repository) Update(ctx context.Context, id uuid.UUID, name, description *string) error {
	// Get current module to preserve values if not updating
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	params := store.UpdateModuleParams{
		ID:          id,
		Name:        sql.NullString{String: current.Name, Valid: true},
		Description: sql.NullString{String: current.Description, Valid: current.Description != ""},
	}

	if name != nil {
		params.Name = sql.NullString{String: *name, Valid: true}
	}

	if description != nil {
		if *description != "" {
			params.Description = sql.NullString{String: *description, Valid: true}
		} else {
			params.Description = sql.NullString{String: "", Valid: false}
		}
	}

	return r.queries.UpdateModule(ctx, params)
}

// Delete deletes a module
func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteModule(ctx, id)
}
