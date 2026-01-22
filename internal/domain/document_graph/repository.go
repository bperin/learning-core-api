package document_graph

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) (*Repository, error) {
	if db == nil {
		return nil, fmt.Errorf("db is required")
	}
	return &Repository{db: db}, nil
}

func (r *Repository) ReplaceGraph(ctx context.Context, documentID uuid.UUID, nodes []Node, edges []Edge) error {
	if documentID == uuid.Nil {
		return fmt.Errorf("document id is required")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM document_graph_edges WHERE document_id = $1", documentID); err != nil {
		return fmt.Errorf("failed to clear edges: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM document_graph_nodes WHERE document_id = $1", documentID); err != nil {
		return fmt.Errorf("failed to clear nodes: %w", err)
	}

	if len(nodes) > 0 {
		stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO document_graph_nodes (id, document_id, node_type, text_content, page_number, metadata)
			VALUES ($1, $2, $3, $4, $5, $6)
		`)
		if err != nil {
			return fmt.Errorf("failed to prepare node insert: %w", err)
		}
		defer stmt.Close()

		for _, node := range nodes {
			var pageNumber sql.NullInt32
			if node.PageNumber != nil {
				pageNumber = sql.NullInt32{Int32: int32(*node.PageNumber), Valid: true}
			}
			metadata := json.RawMessage(nil)
			if len(node.Metadata) > 0 {
				metadata = node.Metadata
			}
			if _, err := stmt.ExecContext(ctx, node.ID, node.DocumentID, node.NodeType, node.Text, pageNumber, metadata); err != nil {
				return fmt.Errorf("failed to insert node: %w", err)
			}
		}
	}

	if len(edges) > 0 {
		stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO document_graph_edges (id, document_id, from_node_id, to_node_id, relation, metadata)
			VALUES ($1, $2, $3, $4, $5, $6)
		`)
		if err != nil {
			return fmt.Errorf("failed to prepare edge insert: %w", err)
		}
		defer stmt.Close()

		for _, edge := range edges {
			metadata := json.RawMessage(nil)
			if len(edge.Metadata) > 0 {
				metadata = edge.Metadata
			}
			if _, err := stmt.ExecContext(ctx, edge.ID, edge.DocumentID, edge.FromNodeID, edge.ToNodeID, edge.Relation, metadata); err != nil {
				return fmt.Errorf("failed to insert edge: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit graph insert: %w", err)
	}

	return nil
}

func (r *Repository) SearchNodes(ctx context.Context, documentID uuid.UUID, query string, limit int) ([]Node, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, document_id, node_type, text_content, page_number, metadata, created_at
		FROM document_graph_nodes
		WHERE document_id = $1 AND text_content ILIKE $2
		ORDER BY created_at DESC
		LIMIT $3
	`, documentID, "%"+query+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search nodes: %w", err)
	}
	defer rows.Close()

	var nodes []Node
	for rows.Next() {
		var node Node
		var pageNumber sql.NullInt32
		var metadata json.RawMessage
		if err := rows.Scan(&node.ID, &node.DocumentID, &node.NodeType, &node.Text, &pageNumber, &metadata, &node.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}
		if pageNumber.Valid {
			value := int(pageNumber.Int32)
			node.PageNumber = &value
		}
		node.Metadata = metadata
		nodes = append(nodes, node)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate node rows: %w", err)
	}

	return nodes, nil
}

func (r *Repository) FetchNeighbors(ctx context.Context, documentID uuid.UUID, nodeIDs []uuid.UUID, limit int) ([]Node, []Edge, error) {
	if len(nodeIDs) == 0 {
		return nil, nil, nil
	}
	if limit <= 0 {
		limit = 20
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, document_id, from_node_id, to_node_id, relation, metadata, created_at
		FROM document_graph_edges
		WHERE document_id = $1 AND (from_node_id = ANY($2) OR to_node_id = ANY($2))
		LIMIT $3
	`, documentID, pq.Array(nodeIDs), limit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch edges: %w", err)
	}
	defer rows.Close()

	var edges []Edge
	neighborIDs := make(map[uuid.UUID]struct{})
	for rows.Next() {
		var edge Edge
		var metadata json.RawMessage
		if err := rows.Scan(&edge.ID, &edge.DocumentID, &edge.FromNodeID, &edge.ToNodeID, &edge.Relation, &metadata, &edge.CreatedAt); err != nil {
			return nil, nil, fmt.Errorf("failed to scan edge: %w", err)
		}
		edge.Metadata = metadata
		edges = append(edges, edge)
		neighborIDs[edge.FromNodeID] = struct{}{}
		neighborIDs[edge.ToNodeID] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to iterate edge rows: %w", err)
	}

	if len(neighborIDs) == 0 {
		return nil, edges, nil
	}

	ids := make([]uuid.UUID, 0, len(neighborIDs))
	for id := range neighborIDs {
		ids = append(ids, id)
	}

	nodeRows, err := r.db.QueryContext(ctx, `
		SELECT id, document_id, node_type, text_content, page_number, metadata, created_at
		FROM document_graph_nodes
		WHERE document_id = $1 AND id = ANY($2)
	`, documentID, pq.Array(ids))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch neighbor nodes: %w", err)
	}
	defer nodeRows.Close()

	var nodes []Node
	for nodeRows.Next() {
		var node Node
		var pageNumber sql.NullInt32
		var metadata json.RawMessage
		if err := nodeRows.Scan(&node.ID, &node.DocumentID, &node.NodeType, &node.Text, &pageNumber, &metadata, &node.CreatedAt); err != nil {
			return nil, nil, fmt.Errorf("failed to scan neighbor node: %w", err)
		}
		if pageNumber.Valid {
			value := int(pageNumber.Int32)
			node.PageNumber = &value
		}
		node.Metadata = metadata
		nodes = append(nodes, node)
	}
	if err := nodeRows.Err(); err != nil {
		return nil, nil, fmt.Errorf("failed to iterate neighbor nodes: %w", err)
	}

	return nodes, edges, nil
}
