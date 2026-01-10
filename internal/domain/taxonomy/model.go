package taxonomy

import (
	"time"

	"github.com/google/uuid"
)

// TaxonomyState represents the lifecycle state for taxonomy nodes and links.
type TaxonomyState string

const (
	TaxonomyStateAIGenerated TaxonomyState = "ai_generated"
	TaxonomyStateApproved    TaxonomyState = "approved"
	TaxonomyStateRejected    TaxonomyState = "rejected"
)

// TaxonomyNode represents a node in the taxonomy tree.
type TaxonomyNode struct {
	ID               uuid.UUID  `json:"id"`
	Name             string     `json:"name"`
	Description      *string    `json:"description,omitempty"`
	ParentID         *uuid.UUID `json:"parent_id,omitempty"`
	Path             string     `json:"path"`
	Depth            int32      `json:"depth"`
	State            string     `json:"state"`
	Confidence       *float64   `json:"confidence,omitempty"`
	SourceDocumentID *uuid.UUID `json:"source_document_id,omitempty"`
	Version          int32      `json:"version"`
	IsActive         bool       `json:"is_active"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
	ApprovedBy       *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// CreateTaxonomyNodeRequest represents the data needed to create a taxonomy node.
type CreateTaxonomyNodeRequest struct {
	Name             string     `json:"name"`
	Description      *string    `json:"description,omitempty"`
	ParentID         *uuid.UUID `json:"parent_id,omitempty"`
	Path             string     `json:"path"`
	Depth            int32      `json:"depth"`
	State            string     `json:"state"`
	Confidence       *float64   `json:"confidence,omitempty"`
	SourceDocumentID *uuid.UUID `json:"source_document_id,omitempty"`
	IsActive         bool       `json:"is_active"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty"`
	ApprovedBy       *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
}

// DocumentTaxonomyLink represents a link between a document and a taxonomy node.
type DocumentTaxonomyLink struct {
	DocumentID     uuid.UUID  `json:"document_id"`
	TaxonomyNodeID uuid.UUID  `json:"taxonomy_node_id"`
	Confidence     *float64   `json:"confidence,omitempty"`
	State          string     `json:"state"`
	CreatedAt      time.Time  `json:"created_at"`
	ApprovedBy     *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt     *time.Time `json:"approved_at,omitempty"`
}

// CreateDocumentTaxonomyLinkRequest represents the data needed to create a document taxonomy link.
type CreateDocumentTaxonomyLinkRequest struct {
	DocumentID     uuid.UUID  `json:"document_id"`
	TaxonomyNodeID uuid.UUID  `json:"taxonomy_node_id"`
	Confidence     *float64   `json:"confidence,omitempty"`
	State          string     `json:"state"`
	ApprovedBy     *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt     *time.Time `json:"approved_at,omitempty"`
}

// UpdateDocumentTaxonomyLinkStateRequest represents the data needed to update link state.
type UpdateDocumentTaxonomyLinkStateRequest struct {
	DocumentID     uuid.UUID  `json:"document_id"`
	TaxonomyNodeID uuid.UUID  `json:"taxonomy_node_id"`
	State          string     `json:"state"`
	ApprovedBy     *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt     *time.Time `json:"approved_at,omitempty"`
}
