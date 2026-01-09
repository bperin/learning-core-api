package synthetic

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func decodeIntents(raw json.RawMessage) []Intent {
	var payload struct {
		Intents []Intent `json:"intents"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	return payload.Intents
}

func decodePlan(raw json.RawMessage) (Plan, error) {
	var plan Plan
	if err := json.Unmarshal(raw, &plan); err != nil {
		return Plan{}, err
	}
	if plan.Title == "" {
		return plan, fmt.Errorf("plan title is required")
	}
	return plan, nil
}

func decodeEvalPayload(raw json.RawMessage) (EvalPayload, error) {
	var payload EvalPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return EvalPayload{}, err
	}
	if payload.Title == "" {
		return payload, fmt.Errorf("eval title is required")
	}
	return payload, nil
}

func buildArtifact(
	artifactType string,
	output json.RawMessage,
	doc DocumentReference,
	prompt PromptTemplate,
	schema SchemaTemplate,
	promptRender string,
) *Artifact {
	createdAt := time.Now().UTC()
	status := "READY"
	model := prompt.Model
	promptText := prompt.Template

	return &Artifact{
		Type:             artifactType,
		Status:           status,
		ReviewerID:       doc.ReviewerID,
		OutputJSON:       output,
		Model:            optionalString(model),
		Prompt:           optionalString(promptText),
		PromptTemplateID: optionalUUID(prompt.ID),
		SchemaTemplateID: optionalUUID(schema.ID),
		ModelParams:      prompt.ModelParams,
		PromptRender:     optionalString(promptRender),
		CreatedAt:        createdAt,
	}
}

func renderPrompt(engine SyntheticEngine, prompt PromptTemplate, inputVars map[string]any) (string, error) {
	if renderer, ok := engine.(PromptRenderer); ok {
		return renderer.RenderPrompt(prompt, inputVars)
	}
	return prompt.Template, nil
}

func optionalUUID(id uuid.UUID) *uuid.UUID {
	if id == uuid.Nil {
		return nil
	}
	return &id
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
