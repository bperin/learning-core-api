package suites

import (
	"time"

	"github.com/google/uuid"
)

// Suite represents an evaluation suite
type Suite struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateSuiteRequest represents the request to create an evaluation suite
type CreateSuiteRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
