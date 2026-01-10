package synthetic

import (
	"context"
)

type SyntheticIntentGenerator struct {
	engine SyntheticEngine
}

func NewSyntheticIntentGenerator(engine SyntheticEngine) *SyntheticIntentGenerator {
	return &SyntheticIntentGenerator{engine: engine}
}

func (g *SyntheticIntentGenerator) GenerateIntents(
	ctx context.Context,
	doc DocumentReference,
	prompt PromptTemplate,
	schema SchemaTemplate,
) ([]Intent, *Artifact, error) {

	inputVars := map[string]any{
		"title": doc.HumanTitle,
	}

	rawJSON, err := g.engine.Generate(ctx, prompt, schema, inputVars)
	if err != nil {
		return nil, nil, err
	}

	intents := decodeIntents(rawJSON)
	promptRender, _ := renderPrompt(g.engine, prompt, inputVars)
	artifact := buildArtifact("INTENTS", rawJSON, doc, prompt, schema, promptRender)

	return intents, artifact, nil
}
