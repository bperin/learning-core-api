package document_graph

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Node struct {
	ID         uuid.UUID       `json:"id"`
	DocumentID uuid.UUID       `json:"document_id"`
	NodeType   string          `json:"node_type"`
	Text       string          `json:"text"`
	PageNumber *int            `json:"page_number,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

type Edge struct {
	ID         uuid.UUID       `json:"id"`
	DocumentID uuid.UUID       `json:"document_id"`
	FromNodeID uuid.UUID       `json:"from_node_id"`
	ToNodeID   uuid.UUID       `json:"to_node_id"`
	Relation   string          `json:"relation"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

type BuildResult struct {
	DocumentID   uuid.UUID `json:"document_id"`
	NodesCreated int       `json:"nodes_created"`
	EdgesCreated int       `json:"edges_created"`
}

type QueryRequest struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

type QueryResponse struct {
	DocumentID uuid.UUID `json:"document_id"`
	Nodes      []Node    `json:"nodes"`
	Edges      []Edge    `json:"edges"`
}
