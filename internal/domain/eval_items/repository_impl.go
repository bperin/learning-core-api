package eval_items

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"

	"learning-core-api/internal/persistance/store"
)

// RepositoryImpl implements the Repository interface using SQLC
type RepositoryImpl struct {
	queries *store.Queries
}

// NewRepository creates a new eval items repository
func NewRepository(queries *store.Queries) Repository {
	return &RepositoryImpl{
		queries: queries,
	}
}

// Create creates a new evaluation item
func (r *RepositoryImpl) Create(ctx context.Context, req *CreateEvalItemRequest) (*EvalItem, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Convert metadata to JSON
	var metadata pqtype.NullRawMessage
	if req.Metadata != nil {
		jsonData, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadata = pqtype.NullRawMessage{RawMessage: jsonData, Valid: true}
	}

	storeItem, err := r.queries.CreateEvalItem(ctx, store.CreateEvalItemParams{
		EvalID:      req.EvalID,
		Prompt:      req.Prompt,
		Options:     req.Options,
		CorrectIdx:  req.CorrectIdx,
		Hint:        toNullString(req.Hint),
		Explanation: toNullString(req.Explanation),
		Metadata:    metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create eval item: %w", err)
	}

	return toDomainEvalItem(&storeItem), nil
}

// GetByID retrieves an evaluation item by ID
func (r *RepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*EvalItem, error) {
	storeItem, err := r.queries.GetEvalItem(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewEvalItemNotFoundError(id)
		}
		return nil, fmt.Errorf("failed to get eval item: %w", err)
	}

	return toDomainEvalItem(&storeItem), nil
}

// List lists evaluation items with pagination
func (r *RepositoryImpl) List(ctx context.Context, req *ListEvalItemsRequest) ([]*EvalItem, error) {
	// If EvalID is specified, use GetEvalItemsByEval instead
	if req.EvalID != nil {
		storeItems, err := r.queries.GetEvalItemsByEval(ctx, *req.EvalID)
		if err != nil {
			return nil, fmt.Errorf("failed to list eval items by eval ID: %w", err)
		}

		// Apply pagination manually since GetEvalItemsByEval doesn't support it
		start := int(req.Offset)
		end := start + int(req.Limit)
		if start >= len(storeItems) {
			return []*EvalItem{}, nil
		}
		if end > len(storeItems) {
			end = len(storeItems)
		}

		items := make([]*EvalItem, end-start)
		for i, storeItem := range storeItems[start:end] {
			items[i] = toDomainEvalItem(&storeItem)
		}
		return items, nil
	}

	// Use general list method
	params := store.ListEvalItemsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	storeItems, err := r.queries.ListEvalItems(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to list eval items: %w", err)
	}

	items := make([]*EvalItem, len(storeItems))
	for i, storeItem := range storeItems {
		items[i] = toDomainEvalItem(&storeItem)
	}

	return items, nil
}

// GetByEvalID retrieves all evaluation items for a specific evaluation
func (r *RepositoryImpl) GetByEvalID(ctx context.Context, evalID uuid.UUID) ([]*EvalItem, error) {
	storeItems, err := r.queries.GetEvalItemsByEval(ctx, evalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get eval items by eval ID: %w", err)
	}

	items := make([]*EvalItem, len(storeItems))
	for i, storeItem := range storeItems {
		items[i] = toDomainEvalItem(&storeItem)
	}

	return items, nil
}

// Search searches evaluation items by prompt text
func (r *RepositoryImpl) Search(ctx context.Context, req *SearchEvalItemsRequest) ([]*EvalItem, error) {
	params := store.SearchEvalItemsByPromptParams{
		Column1: sql.NullString{String: req.Query, Valid: true},
		Limit:   req.Limit,
		Offset:  req.Offset,
	}

	storeItems, err := r.queries.SearchEvalItemsByPrompt(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search eval items: %w", err)
	}

	items := make([]*EvalItem, len(storeItems))
	for i, storeItem := range storeItems {
		items[i] = toDomainEvalItem(&storeItem)
	}

	return items, nil
}

// GetRandom retrieves random evaluation items for a specific evaluation
func (r *RepositoryImpl) GetRandom(ctx context.Context, evalID uuid.UUID, limit int32) ([]*EvalItem, error) {
	params := store.GetRandomEvalItemsParams{
		EvalID: evalID,
		Limit:  limit,
	}

	storeItems, err := r.queries.GetRandomEvalItems(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get random eval items: %w", err)
	}

	items := make([]*EvalItem, len(storeItems))
	for i, storeItem := range storeItems {
		items[i] = toDomainEvalItem(&storeItem)
	}

	return items, nil
}

// Count returns the total number of evaluation items
func (r *RepositoryImpl) Count(ctx context.Context) (int64, error) {
	// Since there's no direct count method, we'll use a list with high limit
	// In a real implementation, you'd want to add a COUNT query to SQLC
	items, err := r.queries.ListEvalItems(ctx, store.ListEvalItemsParams{
		Limit:  10000, // High limit to get all items
		Offset: 0,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to count eval items: %w", err)
	}
	return int64(len(items)), nil
}

// CountByEvalID returns the number of evaluation items for a specific evaluation
func (r *RepositoryImpl) CountByEvalID(ctx context.Context, evalID uuid.UUID) (int64, error) {
	items, err := r.queries.GetEvalItemsByEval(ctx, evalID)
	if err != nil {
		return 0, fmt.Errorf("failed to count eval items by eval ID: %w", err)
	}
	return int64(len(items)), nil
}

// GetWithReviews retrieves an evaluation item with its reviews
func (r *RepositoryImpl) GetWithReviews(ctx context.Context, id uuid.UUID) (*EvalItemWithReviews, error) {
	result, err := r.queries.GetEvalItemWithReviews(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewEvalItemNotFoundError(id)
		}
		return nil, fmt.Errorf("failed to get eval item with reviews: %w", err)
	}

	evalItem := toDomainEvalItem(&store.EvalItem{
		ID:          result.ID,
		EvalID:      result.EvalID,
		Prompt:      result.Prompt,
		Options:     result.Options,
		CorrectIdx:  result.CorrectIdx,
		Hint:        result.Hint,
		Explanation: result.Explanation,
		Metadata:    result.Metadata,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	})

	// TODO: Parse reviews from the result
	// This would require additional SQLC queries or JSON aggregation
	reviews := []*EvalItemReview{}

	return &EvalItemWithReviews{
		EvalItem: *evalItem,
		Reviews:  reviews,
	}, nil
}

// GetWithAnswerStats retrieves evaluation items with answer statistics
func (r *RepositoryImpl) GetWithAnswerStats(ctx context.Context, evalID uuid.UUID) ([]*EvalItemWithAnswerStats, error) {
	results, err := r.queries.GetEvalItemsWithAnswerStats(ctx, evalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get eval items with answer stats: %w", err)
	}

	items := make([]*EvalItemWithAnswerStats, len(results))
	for i, result := range results {
		evalItem := toDomainEvalItem(&store.EvalItem{
			ID:          result.ID,
			EvalID:      result.EvalID,
			Prompt:      result.Prompt,
			Options:     result.Options,
			CorrectIdx:  result.CorrectIdx,
			Hint:        result.Hint,
			Explanation: result.Explanation,
			Metadata:    result.Metadata,
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
		})

		items[i] = &EvalItemWithAnswerStats{
			EvalItem:       *evalItem,
			TotalAnswers:   int32(result.TotalAnswers),
			CorrectAnswers: int32(result.CorrectAnswers),
			SuccessRate:    float64(result.SuccessRate),
		}
	}

	return items, nil
}

// Helper functions

// toDomainEvalItem converts a store.EvalItem to domain.EvalItem
func toDomainEvalItem(storeItem *store.EvalItem) *EvalItem {
	item := &EvalItem{
		ID:         storeItem.ID,
		EvalID:     storeItem.EvalID,
		Prompt:     storeItem.Prompt,
		Options:    storeItem.Options,
		CorrectIdx: storeItem.CorrectIdx,
		CreatedAt:  storeItem.CreatedAt,
		UpdatedAt:  storeItem.UpdatedAt,
	}

	if storeItem.Hint.Valid {
		item.Hint = &storeItem.Hint.String
	}

	if storeItem.Explanation.Valid {
		item.Explanation = &storeItem.Explanation.String
	}

	if storeItem.Metadata.Valid {
		var metadata map[string]interface{}
		if err := json.Unmarshal(storeItem.Metadata.RawMessage, &metadata); err == nil {
			item.Metadata = metadata
		}
	}

	return item
}

// toNullString converts a string pointer to sql.NullString
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}
