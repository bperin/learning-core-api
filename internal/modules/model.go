package modules

import (
	"time"

	"github.com/google/uuid"
)

// Module represents a learning module
type Module struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Title       string    `json:"title"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateModuleRequest represents the request to create a module
type CreateModuleRequest struct {
	TenantID    uuid.UUID `json:"tenant_id"`
	Title       string    `json:"title"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}
