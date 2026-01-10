-- name: CreateTaxonomyNode :one
WITH inserted AS (
  INSERT INTO taxonomy_nodes (
    name, description, parent_id, path, depth, state, confidence,
    source_document_id, version, is_active, created_by, approved_by, approved_at
  ) VALUES (
    $1, $2, $3, $4, $5, $6, $7,
    $8,
    (SELECT COALESCE(MAX(version), 0) + 1 FROM taxonomy_nodes WHERE path = $4),
    $9, $10, $11, $12
  )
  RETURNING *
),
deactivated AS (
  UPDATE taxonomy_nodes SET
    is_active = false
  WHERE path = (SELECT path FROM inserted)
    AND taxonomy_nodes.id != (SELECT id FROM inserted)
    AND (SELECT is_active FROM inserted) = true
)
SELECT * FROM inserted;

-- name: GetTaxonomyNode :one
SELECT * FROM taxonomy_nodes WHERE id = $1 LIMIT 1;

-- name: GetActiveTaxonomyNodeByPath :one
SELECT * FROM taxonomy_nodes WHERE path = $1 AND is_active = true LIMIT 1;

-- name: ListTaxonomyNodesByPrefix :many
SELECT * FROM taxonomy_nodes WHERE path LIKE $1 || '%' ORDER BY path ASC, version DESC;

-- name: ActivateTaxonomyNode :one
WITH target AS (
  SELECT path FROM taxonomy_nodes WHERE taxonomy_nodes.id = $1
),
deactivated AS (
  UPDATE taxonomy_nodes SET
    is_active = false
  WHERE path = (SELECT path FROM target) AND taxonomy_nodes.id != $1
),
activated AS (
  UPDATE taxonomy_nodes SET
    is_active = true
  WHERE taxonomy_nodes.id = $1
  RETURNING *
)
SELECT * FROM activated;

-- name: CreateDocumentTaxonomyLink :one
INSERT INTO document_taxonomy_links (
  document_id, taxonomy_node_id, confidence, state, approved_by, approved_at
) VALUES (
  $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateDocumentTaxonomyLinkState :one
UPDATE document_taxonomy_links SET
  state = $3,
  approved_by = $4,
  approved_at = $5
WHERE document_id = $1 AND taxonomy_node_id = $2
RETURNING *;

-- name: ListDocumentsByTaxonomyPrefix :many
SELECT DISTINCT d.*
FROM documents d
JOIN document_taxonomy_links dtl ON d.id = dtl.document_id
JOIN taxonomy_nodes tn ON tn.id = dtl.taxonomy_node_id
WHERE tn.path LIKE $1 || '%'
  AND tn.state = 'approved'
  AND dtl.state = 'approved'
  AND tn.is_active = true
ORDER BY d.created_at DESC;
