package synthetic

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyntheticPipeline(t *testing.T) {
	engine := NewGenericSyntheticEngine(42)
	intentGen := NewSyntheticIntentGenerator(engine)
	planGen := NewSyntheticPlanGenerator(engine)
	evalGen := NewSyntheticEvalGenerator(engine)

	docID := uuid.New()
	userID := uuid.New()
	doc := DocumentReference{
		ID:          docID,
		SubjectName: "Math",
		Curriculum:  "AP Calc",
		HumanTitle:  "Derivatives",
		TopicTags:   []string{"derivatives", "limits"},
		ReviewerID:  &userID,
		RequestedBy: userID,
		RequestedAt: time.Now().UTC(),
	}

	intentSchema := SchemaTemplate{
		ID:         uuid.New(),
		SchemaType: "intent_extraction",
		Version:    1,
		SchemaJSON: json.RawMessage(`{
      "type": "object",
      "properties": {
        "intents": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "title": {"type": "string", "examples": ["Review derivatives"]},
              "description": {"type": "string"}
            }
          }
        }
      }
    }`),
	}

	planSchema := SchemaTemplate{
		ID:         uuid.New(),
		SchemaType: "plan_generation",
		Version:    1,
		SchemaJSON: json.RawMessage(`{
      "type": "object",
      "properties": {
        "title": {"type": "string"},
        "steps": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "title": {"type": "string"},
              "objectives": {
                "type": "array",
                "items": {"type": "string"}
              }
            }
          }
        }
      }
    }`),
	}

	evalSchema := SchemaTemplate{
		ID:         uuid.New(),
		SchemaType: "eval_generation",
		Version:    1,
		SchemaJSON: json.RawMessage(`{
      "type": "object",
      "properties": {
        "title": {"type": "string"},
        "description": {"type": "string"},
        "items": {
          "type": "array",
          "items": {
            "type": "object",
            "properties": {
              "prompt": {"type": "string"},
              "options": {"type": "array", "minItems": 2, "maxItems": 4, "items": {"type": "string"}},
              "correct_index": {"type": "integer", "minimum": 0, "maximum": 3},
              "hint": {"type": "string"},
              "explanation": {"type": "string"}
            }
          }
        }
      }
    }`),
	}

	intentPrompt := PromptTemplate{ID: uuid.New(), Key: "intent", Version: 1, Template: "Intent for {{.subject}}"}
	planPrompt := PromptTemplate{ID: uuid.New(), Key: "plan", Version: 1, Template: "Plan for {{.title}}"}
	evalPrompt := PromptTemplate{ID: uuid.New(), Key: "eval", Version: 1, Template: "Eval for {{.title}}"}

	t.Logf("document reference: %+v", doc)
	t.Logf("intent prompt: %+v", intentPrompt)
	t.Logf("plan prompt: %+v", planPrompt)
	t.Logf("eval prompt: %+v", evalPrompt)

	ctx := context.Background()
	intentInputs := map[string]any{
		"subject":    doc.SubjectName,
		"curriculum": doc.Curriculum,
		"title":      doc.HumanTitle,
		"tags":       doc.TopicTags,
	}
	logJSON(t, "intent inputs", intentInputs)
	intents, intentArtifact, err := intentGen.GenerateIntents(ctx, doc, intentPrompt, intentSchema)
	require.NoError(t, err)
	require.NotEmpty(t, intents)
	require.NotNil(t, intentArtifact)
	assert.Equal(t, "INTENTS", intentArtifact.Type)
	assert.Equal(t, intentPrompt.ID, *intentArtifact.PromptTemplateID)
	assert.Equal(t, intentSchema.ID, *intentArtifact.SchemaTemplateID)
	logJSON(t, "intent output", json.RawMessage(intentArtifact.OutputJSON))
	logArtifact(t, "intent artifact", intentArtifact)
	if intentArtifact.PromptRender != nil {
		t.Logf("intent prompt render: %s", *intentArtifact.PromptRender)
	}
	logJSON(t, "intents", intents)

	planInputs := map[string]any{
		"subject":    doc.SubjectName,
		"curriculum": doc.Curriculum,
		"title":      doc.HumanTitle,
		"tags":       doc.TopicTags,
		"intents":    intents,
	}
	logJSON(t, "plan inputs", planInputs)
	plan, planArtifact, err := planGen.GeneratePlan(ctx, doc, intents, planPrompt, planSchema)
	require.NoError(t, err)
	require.NotNil(t, planArtifact)
	assert.Equal(t, "PLAN", planArtifact.Type)
	assert.NotEmpty(t, plan.Title)
	logJSON(t, "plan output", json.RawMessage(planArtifact.OutputJSON))
	logArtifact(t, "plan artifact", planArtifact)
	if planArtifact.PromptRender != nil {
		t.Logf("plan prompt render: %s", *planArtifact.PromptRender)
	}
	logJSON(t, "plan", plan)

	evalInputs := map[string]any{
		"subject":    doc.SubjectName,
		"curriculum": doc.Curriculum,
		"title":      doc.HumanTitle,
		"tags":       doc.TopicTags,
		"plan":       plan,
	}
	logJSON(t, "eval inputs", evalInputs)
	eval, evalItems, evalArtifact, err := evalGen.GenerateEval(ctx, doc, plan, evalPrompt, evalSchema)
	require.NoError(t, err)
	require.NotNil(t, eval)
	require.NotNil(t, evalArtifact)
	require.NotEmpty(t, evalItems)
	assert.Equal(t, "EVAL", evalArtifact.Type)
	assert.Equal(t, eval.ID, *evalArtifact.EvalID)
	logJSON(t, "eval output", json.RawMessage(evalArtifact.OutputJSON))
	logArtifact(t, "eval artifact", evalArtifact)
	logJSON(t, "eval domain", eval)
	logJSON(t, "eval items", evalItems)
	if evalArtifact.PromptRender != nil {
		t.Logf("eval prompt render: %s", *evalArtifact.PromptRender)
	}

	for _, item := range evalItems {
		assert.Equal(t, eval.ID, item.EvalID)
		assert.NotEmpty(t, item.Prompt)
		assert.GreaterOrEqual(t, len(item.Options), 2)
	}
}

func logJSON(t *testing.T, label string, value any) {
	t.Helper()
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Logf("%s: <error: %v>", label, err)
		return
	}
	t.Logf("%s: %s", label, payload)
}

func logArtifact(t *testing.T, label string, artifact *Artifact) {
	t.Helper()
	if artifact == nil {
		t.Logf("%s: <nil>", label)
		return
	}
	logJSON(t, label, map[string]any{
		"type":                artifact.Type,
		"status":              artifact.Status,
		"reviewer_id":         artifact.ReviewerID,
		"eval_id":             artifact.EvalID,
		"prompt_template_id":  artifact.PromptTemplateID,
		"schema_template_id":  artifact.SchemaTemplateID,
		"model":               artifact.Model,
		"prompt":              artifact.Prompt,
		"prompt_render":       artifact.PromptRender,
		"model_params":        artifact.ModelParams,
		"created_at":          artifact.CreatedAt,
		"output_json_preview": string(artifact.OutputJSON),
	})
}
