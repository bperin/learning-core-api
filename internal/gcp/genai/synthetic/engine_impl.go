package synthetic

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand/v2"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
)

type GenericSyntheticEngine struct {
	rng *rand.Rand
}

func NewGenericSyntheticEngine(seed uint64) *GenericSyntheticEngine {
	return &GenericSyntheticEngine{
		rng: rand.New(rand.NewPCG(seed, seed+1)),
	}
}

func (e *GenericSyntheticEngine) Generate(_ context.Context, prompt PromptTemplate, schema SchemaTemplate, inputVars map[string]any) (json.RawMessage, error) {
	_, err := e.RenderPrompt(prompt, inputVars)
	if err != nil {
		return nil, err
	}

	var schemaDef map[string]any
	if err := json.Unmarshal(schema.SchemaJSON, &schemaDef); err != nil {
		return nil, err
	}

	output := e.walk(schemaDef)
	return json.Marshal(output)
}

func (e *GenericSyntheticEngine) RenderPrompt(prompt PromptTemplate, inputVars map[string]any) (string, error) {
	if strings.TrimSpace(prompt.Template) == "" {
		return "", nil
	}

	tmpl, err := template.New("prompt").Option("missingkey=zero").Parse(prompt.Template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, inputVars); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (e *GenericSyntheticEngine) walk(schema any) any {
	switch typed := schema.(type) {
	case map[string]any:
		return e.walkObjectSchema(typed)
	case []any:
		if len(typed) == 0 {
			return []any{}
		}
		return e.walkArraySchema(map[string]any{"items": typed[0]})
	default:
		return nil
	}
}

func (e *GenericSyntheticEngine) walkObjectSchema(schema map[string]any) any {
	if value, ok := schema["const"]; ok {
		return value
	}

	if enumValues, ok := schema["enum"].([]any); ok && len(enumValues) > 0 {
		return enumValues[e.rng.IntN(len(enumValues))]
	}

	if options := schemaOptions(schema, "anyOf", "oneOf"); len(options) > 0 {
		choice := options[e.rng.IntN(len(options))]
		return e.walk(choice)
	}

	if options := schemaOptions(schema, "allOf"); len(options) > 0 {
		return e.mergeAllOf(options)
	}

	if schemaType, ok := schema["type"].(string); ok {
		switch schemaType {
		case "object":
			return e.walkProperties(schema)
		case "array":
			return e.walkArraySchema(schema)
		case "string":
			return e.fakeString(schema)
		case "integer":
			return e.fakeInteger(schema)
		case "number":
			return e.fakeNumber(schema)
		case "boolean":
			return e.rng.IntN(2) == 0
		}
	}

	if _, ok := schema["properties"].(map[string]any); ok {
		return e.walkProperties(schema)
	}

	return nil
}

func (e *GenericSyntheticEngine) walkProperties(schema map[string]any) map[string]any {
	properties, _ := schema["properties"].(map[string]any)
	keys := sortedKeys(properties)
	output := make(map[string]any, len(keys))
	for _, key := range keys {
		output[key] = e.walk(properties[key])
	}
	return output
}

func (e *GenericSyntheticEngine) walkArraySchema(schema map[string]any) []any {
	itemsSchema, ok := schema["items"]
	if !ok {
		return []any{}
	}

	minItems := int64From(schema["minItems"], 1)
	maxItems := int64From(schema["maxItems"], minItems+2)
	if maxItems < minItems {
		maxItems = minItems
	}
	count := minItems
	if maxItems > minItems {
		count += int64(e.rng.IntN(int(maxItems-minItems) + 1))
	}

	output := make([]any, 0, count)
	for i := int64(0); i < count; i++ {
		output = append(output, e.walk(itemsSchema))
	}
	return output
}

func (e *GenericSyntheticEngine) mergeAllOf(options []any) any {
	merged := make(map[string]any)
	for _, option := range options {
		if option == nil {
			continue
		}
		value := e.walk(option)
		child, ok := value.(map[string]any)
		if !ok {
			return value
		}
		for key, val := range child {
			merged[key] = val
		}
	}
	return merged
}

func (e *GenericSyntheticEngine) fakeString(schema map[string]any) string {
	if examples, ok := schema["examples"].([]any); ok && len(examples) > 0 {
		if choice, ok := examples[e.rng.IntN(len(examples))].(string); ok {
			return choice
		}
	}

	if def, ok := schema["default"].(string); ok {
		return def
	}

	if format, ok := schema["format"].(string); ok {
		switch format {
		case "uuid":
			return uuid.New().String()
		case "date-time":
			return time.Now().UTC().Format(time.RFC3339)
		case "email":
			return "synthetic@example.com"
		}
	}

	minLength := int(int64From(schema["minLength"], 5))
	maxLength := int(int64From(schema["maxLength"], int64(minLength+10)))
	length := minLength
	if maxLength > minLength {
		length += e.rng.IntN(maxLength - minLength + 1)
	}
	return randomWords(e.rng, length)
}

func (e *GenericSyntheticEngine) fakeInteger(schema map[string]any) int64 {
	min := int64From(schema["minimum"], 0)
	max := int64From(schema["maximum"], min+10)
	if max < min {
		max = min
	}
	if max == min {
		return min
	}
	return min + int64(e.rng.IntN(int(max-min)+1))
}

func (e *GenericSyntheticEngine) fakeNumber(schema map[string]any) float64 {
	min := floatFrom(schema["minimum"], 0)
	max := floatFrom(schema["maximum"], min+10)
	if max < min {
		max = min
	}
	if max == min {
		return min
	}
	return min + e.rng.Float64()*(max-min)
}

func schemaOptions(schema map[string]any, key string, fallback ...string) []any {
	if options, ok := schema[key].([]any); ok {
		return options
	}
	for _, next := range fallback {
		if options, ok := schema[next].([]any); ok {
			return options
		}
	}
	return nil
}

func int64From(value any, fallback int64) int64 {
	switch v := value.(type) {
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int64:
		return v
	default:
		return fallback
	}
}

func floatFrom(value any, fallback float64) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	case int64:
		return float64(v)
	default:
		return fallback
	}
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func randomWords(rng *rand.Rand, length int) string {
	if length <= 0 {
		return ""
	}

	letters := []rune("abcdefghijklmnopqrstuvwxyz ")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(letters[rng.IntN(len(letters))])
	}
	return strings.TrimSpace(b.String())
}
