package synthetic

import (
	"context"
	"encoding/json"
)

type SyntheticEngine interface {
	Generate(
		ctx context.Context,
		prompt PromptTemplate,
		schema SchemaTemplate,
		inputVars map[string]any, // anything your pipeline passes
	) (json.RawMessage, error)
}

type PromptRenderer interface {
	RenderPrompt(prompt PromptTemplate, inputVars map[string]any) (string, error)
}
