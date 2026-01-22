-- Document graph storage for structured document retrieval
CREATE TABLE IF NOT EXISTS document_graph_nodes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  node_type TEXT NOT NULL,
  text_content TEXT,
  page_number INT,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS document_graph_edges (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  from_node_id UUID NOT NULL REFERENCES document_graph_nodes(id) ON DELETE CASCADE,
  to_node_id UUID NOT NULL REFERENCES document_graph_nodes(id) ON DELETE CASCADE,
  relation TEXT NOT NULL,
  metadata JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS document_graph_nodes_document_id_idx ON document_graph_nodes(document_id);
CREATE INDEX IF NOT EXISTS document_graph_edges_document_id_idx ON document_graph_edges(document_id);
CREATE INDEX IF NOT EXISTS document_graph_edges_from_node_idx ON document_graph_edges(from_node_id);
CREATE INDEX IF NOT EXISTS document_graph_edges_to_node_idx ON document_graph_edges(to_node_id);
