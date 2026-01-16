package generation

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"learning-core-api/internal/domain/prompt_templates"
	"learning-core-api/internal/domain/schema_templates"
	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/testutil"
	"learning-core-api/internal/utils"
)

func TestQuestionsTemplatesSeeded(t *testing.T) {
	tx, cleanup := testutil.NewTestTx(t)
	defer cleanup()

	ctx := context.Background()
	queries := store.New(tx)
	promptRepo := prompt_templates.NewRepository(queries)
	schemaRepo := schema_templates.NewRepository(queries)
	generationType := utils.GenerationTypeQuestions.String()

	prompt, err := promptRepo.GetActiveByGenerationType(ctx, generationType)
	require.NoError(t, err)
	assert.Equal(t, generationType, prompt.GenerationType)
	assert.NotEmpty(t, prompt.Template)

	schema, err := schemaRepo.GetActiveByGenerationType(ctx, generationType)
	require.NoError(t, err)
	assert.Equal(t, generationType, schema.GenerationType)
	require.NotEmpty(t, schema.SchemaJSON)

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(schema.SchemaJSON, &payload))
	props, ok := payload["properties"].(map[string]interface{})
	require.True(t, ok, "schema should define properties")
	_, ok = props["questions"]
	assert.True(t, ok, "schema should define questions")
}
