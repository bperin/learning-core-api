package taxonomy

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type generatedTaxonomyPayload struct {
	ProposedTaxonomy []generatedTaxonomyNode `json:"proposed_taxonomy"`
}

type generatedTaxonomyNode struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Children    []json.RawMessage `json:"children"`
}

func IngestGeneratedTaxonomy(ctx context.Context, repo Repository, documentID uuid.UUID, createdBy uuid.UUID, payload json.RawMessage) ([]*TaxonomyNode, error) {
	if repo == nil {
		return nil, fmt.Errorf("repository is required")
	}
	if documentID == uuid.Nil {
		return nil, fmt.Errorf("document id is required")
	}
	if len(payload) == 0 {
		return nil, fmt.Errorf("taxonomy payload is required")
	}

	var parsed generatedTaxonomyPayload
	if err := json.Unmarshal(payload, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse taxonomy payload: %w", err)
	}

	var created []*TaxonomyNode
	for _, node := range parsed.ProposedTaxonomy {
		nodes, err := ingestNode(ctx, repo, documentID, createdBy, node, uuid.Nil, "", 0)
		if err != nil {
			return nil, err
		}
		created = append(created, nodes...)
	}
	return created, nil
}

func ingestNode(ctx context.Context, repo Repository, documentID uuid.UUID, createdBy uuid.UUID, node generatedTaxonomyNode, parentID uuid.UUID, parentPath string, depth int32) ([]*TaxonomyNode, error) {
	name := normalizeTaxonomyName(node.Name)
	if name == "" {
		return nil, fmt.Errorf("taxonomy node name is required")
	}

	if isDuplicateChildName(node.Children) {
		return nil, fmt.Errorf("taxonomy node %q has duplicated child keys, possible invalid schema", name)
	}

	if isSkippableTaxonomyContainer(name) {
		var created []*TaxonomyNode
		for _, raw := range node.Children {
			child, err := parseChildNode(raw)
			if err != nil {
				return nil, err
			}
			if child == nil {
				continue
			}
			childNodes, err := ingestNode(ctx, repo, documentID, createdBy, *child, parentID, parentPath, depth)
			if err != nil {
				return nil, err
			}
			created = append(created, childNodes...)
		}
		return created, nil
	}

	path := buildTaxonomyPath(parentPath, normalizeTaxonomyPathSegment(name))
	req := CreateTaxonomyNodeRequest{
		Name:             name,
		Description:      normalizeDescription(node.Description),
		ParentID:         ptrUUID(parentID),
		Path:             path,
		Depth:            depth,
		State:            string(TaxonomyStateAIGenerated),
		IsActive:         true,
		CreatedBy:        &createdBy,
		SourceDocumentID: &documentID,
	}

	createdNode, err := repo.CreateNode(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create taxonomy node %q: %w", path, err)
	}

	_, err = repo.CreateDocumentLink(ctx, CreateDocumentTaxonomyLinkRequest{
		DocumentID:     documentID,
		TaxonomyNodeID: createdNode.ID,
		State:          string(TaxonomyStateAIGenerated),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to link taxonomy node %q to document: %w", path, err)
	}

	created := []*TaxonomyNode{createdNode}
	for _, raw := range node.Children {
		child, err := parseChildNode(raw)
		if err != nil {
			return nil, err
		}
		if child == nil {
			continue
		}
		childNodes, err := ingestNode(ctx, repo, documentID, createdBy, *child, createdNode.ID, path, depth+1)
		if err != nil {
			return nil, err
		}
		created = append(created, childNodes...)
	}

	return created, nil
}

func isDuplicateChildName(children []json.RawMessage) bool {
	for _, raw := range children {
		var asString string
		if err := json.Unmarshal(raw, &asString); err == nil {
			normalized := strings.ToLower(strings.TrimSpace(asString))
			if normalized == "name" || normalized == "description" || normalized == "children" {
				return true
			}
		}
	}
	return false
}

func parseChildNode(raw json.RawMessage) (*generatedTaxonomyNode, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	var asObject generatedTaxonomyNode
	if err := json.Unmarshal(raw, &asObject); err == nil && strings.TrimSpace(asObject.Name) != "" {
		return &asObject, nil
	}

	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		name := normalizeTaxonomyName(asString)
		if name == "" || isReservedTaxonomyToken(name) || isLikelySentence(name) {
			return nil, nil
		}
		return &generatedTaxonomyNode{Name: name}, nil
	}

	return nil, fmt.Errorf("invalid taxonomy child: %s", string(raw))
}

func normalizeTaxonomyName(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ""
	}
	trimmed = strings.ReplaceAll(trimmed, "/", "-")
	fields := strings.Fields(trimmed)
	if len(fields) == 0 {
		return ""
	}
	return strings.Join(fields, " ")
}

func isReservedTaxonomyToken(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "name", "description", "children":
		return true
	default:
		return false
	}
}

func isSkippableTaxonomyContainer(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "taxonomy", "taxonomic groups", "taxonomic group":
		return true
	default:
		return false
	}
}

func isLikelySentence(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return false
	}
	if strings.ContainsAny(trimmed, ".;:") {
		return true
	}
	words := strings.Fields(trimmed)
	if len(words) > 10 {
		return true
	}
	return len(trimmed) > 120
}

func normalizeTaxonomyPathSegment(name string) string {
	lowered := strings.ToLower(strings.TrimSpace(name))
	if lowered == "" {
		return ""
	}
	var b strings.Builder
	prevDash := false
	for _, r := range lowered {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
			prevDash = false
		case r >= '0' && r <= '9':
			b.WriteRune(r)
			prevDash = false
		default:
			if !prevDash {
				b.WriteRune('-')
				prevDash = true
			}
		}
	}
	segment := strings.Trim(b.String(), "-")
	segment = strings.ReplaceAll(segment, "--", "-")
	return segment
}

func normalizeDescription(description string) *string {
	trimmed := strings.TrimSpace(description)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func buildTaxonomyPath(parentPath, name string) string {
	if parentPath == "" {
		return name
	}
	return parentPath + "/" + name
}

func ptrUUID(value uuid.UUID) *uuid.UUID {
	if value == uuid.Nil {
		return nil
	}
	return &value
}
