package synthetic

import "context"

type SyntheticPlanGenerator struct {
	engine SyntheticEngine
}

func NewSyntheticPlanGenerator(engine SyntheticEngine) *SyntheticPlanGenerator {
	return &SyntheticPlanGenerator{engine: engine}
}

func (g *SyntheticPlanGenerator) GeneratePlan(
	ctx context.Context,
	doc DocumentReference,
	intents []Intent,
	prompt PromptTemplate,
	schema SchemaTemplate,
) (Plan, *Artifact, error) {
	inputVars := map[string]any{
		"title":   doc.HumanTitle,
		"intents": intents,
	}

	rawJSON, err := g.engine.Generate(ctx, prompt, schema, inputVars)
	if err != nil {
		return Plan{}, nil, err
	}

	plan, err := decodePlan(rawJSON)
	if err != nil {
		return Plan{}, nil, err
	}

	promptRender, _ := renderPrompt(g.engine, prompt, inputVars)
	artifact := buildArtifact("PLAN", rawJSON, doc, prompt, schema, promptRender)

	return plan, artifact, nil
}
