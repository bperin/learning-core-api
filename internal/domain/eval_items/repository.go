package eval_items

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Repository defines the interface for eval item data operations
type Repository interface {
	// CRUD operations
	Create(ctx context.Context, req *CreateEvalItemRequest) (*EvalItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*EvalItem, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateEvalItemRequest) (*EvalItem, error)
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	List(ctx context.Context, req *ListEvalItemsRequest) ([]*EvalItem, error)
	GetByEvalID(ctx context.Context, evalID uuid.UUID) ([]*EvalItem, error)
	Search(ctx context.Context, req *SearchEvalItemsRequest) ([]*EvalItem, error)
	GetRandom(ctx context.Context, evalID uuid.UUID, limit int32) ([]*EvalItem, error)

	// Aggregate operations
	Count(ctx context.Context) (int64, error)
	CountByEvalID(ctx context.Context, evalID uuid.UUID) (int64, error)

	// Business logic operations
	GetWithReviews(ctx context.Context, id uuid.UUID) (*EvalItemWithReviews, error)
	GetWithAnswerStats(ctx context.Context, evalID uuid.UUID) ([]*EvalItemWithAnswerStats, error)
}

// EvalItemWithReviews represents an eval item with its associated reviews
type EvalItemWithReviews struct {
	EvalItem
	Reviews []*EvalItemReview `json:"reviews"`
}

// EvalItemReview represents a review of an eval item
type EvalItemReview struct {
	ID         uuid.UUID `json:"id"`
	ReviewerID uuid.UUID `json:"reviewer_id"`
	Verdict    string    `json:"verdict"`
	Reasons    []string  `json:"reasons"`
	Comments   *string   `json:"comments,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// EvalItemWithAnswerStats represents an eval item with answer statistics
type EvalItemWithAnswerStats struct {
	EvalItem
	TotalAnswers   int32   `json:"total_answers"`
	CorrectAnswers int32   `json:"correct_answers"`
	SuccessRate    float64 `json:"success_rate"`
}
