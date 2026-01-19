-- +goose Up
-- +goose StatementBegin
-- =====================================================
-- Taxonomy Node Audit Trail
-- =====================================================
-- This migration is a placeholder for future taxonomy node auditing
-- The taxonomy_nodes table may not exist yet, so we skip this migration
-- for now to avoid breaking the build.

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_taxonomy_nodes_updated_at ON taxonomy_nodes;
DROP TRIGGER IF EXISTS audit_taxonomy_nodes_update ON taxonomy_nodes;
DROP FUNCTION IF EXISTS log_taxonomy_node_update();

DROP TABLE IF EXISTS taxonomy_node_updates;

ALTER TABLE taxonomy_nodes
  DROP COLUMN IF EXISTS updated_by;
-- +goose StatementEnd
