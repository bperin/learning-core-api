package synthetic

import (
	"context"
	"time"

	"github.com/google/uuid"

	"learning-core-api/internal/domain/eval_items"
	"learning-core-api/internal/domain/evals"
)

type SyntheticEvalGenerator struct {
	engine SyntheticEngine
}

func NewSyntheticEvalGenerator(engine SyntheticEngine) *SyntheticEvalGenerator {
	return &SyntheticEvalGenerator{engine: engine}
}

func (g *SyntheticEvalGenerator) GenerateEval(
	ctx context.Context,
	doc DocumentReference,
	plan Plan,
	prompt PromptTemplate,
	schema SchemaTemplate,
) (*evals.Eval, []*eval_items.EvalItem, *Artifact, error) {
	inputVars := map[string]any{
		"title": doc.HumanTitle,
		"plan":  plan,
	}

	rawJSON, err := g.engine.Generate(ctx, prompt, schema, inputVars)
	if err != nil {
		return nil, nil, nil, err
	}

	payload, err := decodeEvalPayload(rawJSON)
	if err != nil {
		return nil, nil, nil, err
	}

	createdAt := time.Now().UTC()
	newEvalID := uuid.New()
	status := evals.EvalStatusDraft
	userID := doc.RequestedBy

	newEval := &evals.Eval{
		ID:          newEvalID,
		Title:       payload.Title,
		Status:      status,
		UserID:      userID,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
		Description: optionalEvalDescription(payload.Description),
	}

	items := make([]*eval_items.EvalItem, 0, len(payload.Items))
	for _, item := range payload.Items {
		options := normalizeOptions(item.Options)
		correctIdx := item.CorrectIndex
		if correctIdx < 0 || int(correctIdx) >= len(options) {
			correctIdx = 0
		}
		itemID := uuid.New()
		items = append(items, &eval_items.EvalItem{
			ID:          itemID,
			EvalID:      newEvalID,
			Prompt:      item.Prompt,
			Options:     options,
			CorrectIdx:  correctIdx,
			Hint:        optionalString(item.Hint),
			Explanation: optionalString(item.Explanation),
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		})
	}

	promptRender, _ := renderPrompt(g.engine, prompt, inputVars)
	artifact := buildArtifact("EVAL", rawJSON, doc, prompt, schema, promptRender)
	artifact.EvalID = &newEvalID

	return newEval, items, artifact, nil
}

func normalizeOptions(options []string) []string {
	if len(options) >= 2 {
		return options
	}
	return []string{"Option A", "Option B"}
}

func optionalEvalDescription(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
