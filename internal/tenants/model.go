package tenants

import (
	"time"

	"github.com/google/uuid"
)

// Tenant represents a tenant in the system
type Tenant struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateTenantRequest represents the request to create a tenant
type CreateTenantRequest struct {
	Name string `json:"name"`
}
