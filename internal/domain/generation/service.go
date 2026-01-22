package generation

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/google/uuid"

	"learning-core-api/internal/domain/artifacts"
	"learning-core-api/internal/domain/document_graph"
	"learning-core-api/internal/domain/model_configs"
	"learning-core-api/internal/domain/prompt_templates"
	"learning-core-api/internal/domain/schema_templates"
	"learning-core-api/internal/domain/system_instructions"
	"learning-core-api/internal/persistance/store"
)

type Service struct {
	modelConfigs       model_configs.Repository
	promptTemplates    prompt_templates.Repository
	systemInstructions system_instructions.Repository
	schemaTemplates    schema_templates.Repository
	artifactsService   *artifacts.Service
	generator          Generator
	graphRepo          *document_graph.Repository
}

func NewService(db *sql.DB, artifactsService *artifacts.Service, generator Generator, graphRepo *document_graph.Repository) (*Service, error) {
	if db == nil {
		return nil, fmt.Errorf("db is required")
	}

	queries := store.New(db)
	return &Service{
		modelConfigs:       model_configs.NewRepository(queries),
		promptTemplates:    prompt_templates.NewRepository(queries),
		systemInstructions: system_instructions.NewRepository(queries),
		schemaTemplates:    schema_templates.NewRepository(queries),
		artifactsService:   artifactsService,
		generator:          generator,
		graphRepo:          graphRepo,
	}, nil
}

func (s *Service) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	if req.Instructions.GenerationType != "" && req.Output.GenerationType != "" && req.Instructions.GenerationType != req.Output.GenerationType {
		return nil, fmt.Errorf("generation_type mismatch between instructions and output")
	}

	// 1. Resolve Model Configuration
	resolvedModel, err := s.resolveModelConfig(ctx, req.ModelConfigID)
	if err != nil {
		return nil, err
	}

	// 2. Fetch/Resolve Instructions (Prompt + System Instructions)
	promptText, systemInstr, promptTmplID, err := s.resolveInstructions(ctx, req.Instructions)
	if err != nil {
		return nil, err
	}

	// 3. Fetch/Resolve Output Schema
	responseSchema, schemaTmplID, err := s.resolveOutputConfig(ctx, req.Output)
	if err != nil {
		return nil, err
	}

	// 4. Apply Graph RAG tool context
	tools, promptText, err := s.applyGraphTools(ctx, req, promptText)
	if err != nil {
		return nil, err
	}

	// 5. Call the generator implementation
	resp, err := s.generator.Generate(ctx, GeneratorRequest{
		Prompt:            promptText,
		SystemInstruction: systemInstr,
		OutputSchema:      responseSchema,
		Tools:             tools,
		Model:             resolvedModel,
	})
	if err != nil {
		modelName := modelNameForArtifact(resolvedModel)
		modelParams, meta, metaErr := buildArtifactMetadata(req, resolvedModel, systemInstr)
		if metaErr != nil {
			return nil, metaErr
		}
		s.saveArtifact(ctx, req, promptText, promptTmplID, schemaTmplID, modelName, modelParams, meta, "", nil, err.Error(), nil)
		return nil, fmt.Errorf("genai call failed: %w", err)
	}

	// 5. Extract Output
	outputText := resp.OutputText

	var outputJSON json.RawMessage
	if req.Output.Format == "json" || req.Output.GenerationType != "" || req.Output.InlineSchema != nil {
		outputJSON = json.RawMessage(outputText)
	}

	// 6. Save Artifact
	modelName := modelNameForArtifact(resolvedModel)
	if resp.ModelUsed != "" {
		modelName = resp.ModelUsed
	}

	modelParams, meta, metaErr := buildArtifactMetadata(req, resolvedModel, systemInstr)
	if metaErr != nil {
		return nil, metaErr
	}

	artifactID, saveErr := s.saveArtifact(ctx, req, promptText, promptTmplID, schemaTmplID, modelName, modelParams, meta, outputText, outputJSON, "", resp.GroundingMetadata)
	if saveErr != nil {
		return nil, fmt.Errorf("failed to save artifact: %w", saveErr)
	}

	return &GenerateResponse{
		ArtifactID:        artifactID,
		OutputText:        outputText,
		OutputJSON:        outputJSON,
		FinishReason:      resp.FinishReason,
		ModelUsed:         modelName,
		GroundingMetadata: resp.GroundingMetadata,
	}, nil
}

func (s *Service) resolveModelConfig(ctx context.Context, id uuid.UUID) (*ModelConfig, error) {
	var dbConfig *model_configs.ModelConfig
	var err error

	if id == uuid.Nil {
		dbConfig, err = s.modelConfigs.GetActive(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch active model config: %w", err)
		}
	} else {
		dbConfig, err = s.modelConfigs.GetByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch model config %q: %w", id, err)
		}
	}

	baseConfig := &ModelConfig{
		Name: dbConfig.ModelName,
	}

	temperature := float32(dbConfig.Temperature)
	maxTokens := dbConfig.MaxTokens
	topP := float32(dbConfig.TopP)
	topK := float32(dbConfig.TopK)

	baseConfig.Temperature = &temperature
	baseConfig.MaxTokens = &maxTokens
	baseConfig.TopP = &topP
	baseConfig.TopK = &topK
	baseConfig.MimeType = dbConfig.MimeType

	if baseConfig.Name == "" {
		return nil, fmt.Errorf("resolved model config is incomplete: missing name")
	}

	return baseConfig, nil
}

func (s *Service) resolveInstructions(ctx context.Context, inst Instructions) (string, string, uuid.UUID, error) {
	// 1. Resolve System Instruction
	var systemInstr string
	if inst.SystemInstructionID != nil {
		var sys *system_instructions.SystemInstruction
		var err error
		sys, err = s.systemInstructions.GetByID(ctx, *inst.SystemInstructionID)

		if err != nil {
			return "", "", uuid.Nil, fmt.Errorf("failed to fetch system instruction %q: %w", *inst.SystemInstructionID, err)
		}
		systemInstr = sys.Text
	}

	if inst.Inline != "" {
		return inst.Inline, systemInstr, uuid.Nil, nil
	}
	if inst.GenerationType == "" {
		return "", "", uuid.Nil, fmt.Errorf("generation type is required")
	}

	var promptTmpl *prompt_templates.PromptTemplate
	var err error
	if inst.PromptVersion > 0 {
		promptTmpl, err = s.promptTemplates.GetByGenerationTypeAndVersion(ctx, inst.GenerationType, inst.PromptVersion)
	} else {
		promptTmpl, err = s.promptTemplates.GetActiveByGenerationType(ctx, inst.GenerationType)
	}

	if err != nil {
		return "", "", uuid.Nil, fmt.Errorf("failed to fetch prompt template %q: %w", inst.GenerationType, err)
	}

	tmpl, err := template.New("prompt").Parse(promptTmpl.Template)
	if err != nil {
		return "", "", promptTmpl.ID, fmt.Errorf("failed to parse prompt template: %w", err)
	}

	var renderedPrompt bytes.Buffer
	if err := tmpl.Execute(&renderedPrompt, inst.Variables); err != nil {
		return "", "", promptTmpl.ID, fmt.Errorf("failed to render prompt: %w", err)
	}

	return renderedPrompt.String(), systemInstr, promptTmpl.ID, nil
}

func (s *Service) resolveOutputConfig(ctx context.Context, out OutputConfig) (json.RawMessage, uuid.UUID, error) {
	if out.InlineSchema != nil {
		var parsed interface{}
		if err := json.Unmarshal(out.InlineSchema, &parsed); err != nil {
			return nil, uuid.Nil, fmt.Errorf("failed to parse inline schema: %w", err)
		}
		return out.InlineSchema, uuid.Nil, nil
	}

	if out.GenerationType == "" {
		return nil, uuid.Nil, nil
	}

	var schemaTmpl *schema_templates.SchemaTemplate
	var err error
	if out.SchemaVersion > 0 {
		return nil, uuid.Nil, fmt.Errorf("schema_version is not supported; use the active schema")
	}
	schemaTmpl, err = s.schemaTemplates.GetActiveByGenerationType(ctx, out.GenerationType)

	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("failed to fetch schema template %q: %w", out.GenerationType, err)
	}

	return schemaTmpl.SchemaJSON, schemaTmpl.ID, nil
}

type graphToolConfig struct {
	Query string `json:"query"`
	Limit int    `json:"limit,omitempty"`
}

func (s *Service) applyGraphTools(ctx context.Context, req GenerateRequest, prompt string) ([]ToolConfig, string, error) {
	if len(req.Tools) == 0 {
		return req.Tools, prompt, nil
	}

	filtered := make([]ToolConfig, 0, len(req.Tools))
	graphContext := ""

	for _, tool := range req.Tools {
		if tool.Type != "graph_rag" {
			filtered = append(filtered, tool)
			continue
		}

		if s.graphRepo == nil {
			return nil, "", fmt.Errorf("graph repository is required for graph_rag")
		}
		if req.Target.DocumentID == nil {
			return nil, "", fmt.Errorf("document_id is required for graph_rag")
		}

		var cfg graphToolConfig
		if len(tool.Config) > 0 {
			if err := json.Unmarshal(tool.Config, &cfg); err != nil {
				return nil, "", fmt.Errorf("failed to parse graph_rag config: %w", err)
			}
		}

		query := strings.TrimSpace(cfg.Query)
		if query == "" {
			query = prompt
		}
		if query == "" {
			return nil, "", fmt.Errorf("graph_rag query is required")
		}

		limit := cfg.Limit
		if limit <= 0 {
			limit = 8
		}

		matched, err := s.graphRepo.SearchNodes(ctx, *req.Target.DocumentID, query, limit)
		if err != nil {
			return nil, "", err
		}

		baseIDs := make([]uuid.UUID, 0, len(matched))
		for _, node := range matched {
			baseIDs = append(baseIDs, node.ID)
		}

		neighborNodes, edges, err := s.graphRepo.FetchNeighbors(ctx, *req.Target.DocumentID, baseIDs, limit*4)
		if err != nil {
			return nil, "", err
		}

		nodeMap := map[uuid.UUID]document_graph.Node{}
		for _, node := range matched {
			nodeMap[node.ID] = node
		}
		for _, node := range neighborNodes {
			nodeMap[node.ID] = node
		}

		graphContext = buildGraphContext(nodeMap, edges)
	}

	if graphContext == "" {
		return filtered, prompt, nil
	}

	augmented := fmt.Sprintf("%s\n\n[Graph Context]\n%s", prompt, graphContext)
	return filtered, augmented, nil
}

func buildGraphContext(nodes map[uuid.UUID]document_graph.Node, edges []document_graph.Edge) string {
	if len(nodes) == 0 {
		return ""
	}

	var builder strings.Builder
	for _, node := range nodes {
		label := truncate(node.Text, 240)
		if label == "" {
			label = node.NodeType
		}
		builder.WriteString(fmt.Sprintf("%s: %s\n", node.NodeType, label))
	}

	if len(edges) > 0 {
		builder.WriteString("Relations:\n")
		for _, edge := range edges {
			fromNode, okFrom := nodes[edge.FromNodeID]
			toNode, okTo := nodes[edge.ToNodeID]
			if !okFrom || !okTo {
				continue
			}
			fromLabel := truncate(fromNode.Text, 120)
			toLabel := truncate(toNode.Text, 120)
			builder.WriteString(fmt.Sprintf("- %s -> %s (%s)\n", fromLabel, toLabel, edge.Relation))
		}
	}

	return strings.TrimSpace(builder.String())
}

func truncate(text string, max int) string {
	value := strings.TrimSpace(text)
	if max <= 0 || len(value) <= max {
		return value
	}
	if max < 4 {
		return value[:max]
	}
	return value[:max-3] + "..."
}

func (s *Service) saveArtifact(ctx context.Context, req GenerateRequest, promptText string, promptTmplID, schemaTmplID uuid.UUID, modelName string, modelParams json.RawMessage, meta json.RawMessage, outputText string, outputJSON json.RawMessage, errorMsg string, groundingMetadata json.RawMessage) (uuid.UUID, error) {
	status := "READY"
	if errorMsg != "" {
		status = "ERROR"
	}

	generationType := req.Instructions.GenerationType
	if generationType == "" {
		generationType = req.Output.GenerationType
	}

	// Merge grounding metadata into meta if available
	if len(groundingMetadata) > 0 && len(meta) > 0 {
		var metaObj map[string]interface{}
		var groundingObj map[string]interface{}
		if err := json.Unmarshal(meta, &metaObj); err == nil {
			if err := json.Unmarshal(groundingMetadata, &groundingObj); err == nil {
				metaObj["grounding"] = groundingObj
				if merged, err := json.Marshal(metaObj); err == nil {
					meta = merged
				}
			}
		}
	} else if len(groundingMetadata) > 0 {
		meta = groundingMetadata
	}

	params := artifacts.CreateArtifactParams{
		Type:             "OTHER",
		GenerationType:   generationType,
		Status:           status,
		UserID:           req.UserID,
		DocumentID:       uuid.NullUUID{UUID: ptrToUUID(req.Target.DocumentID), Valid: req.Target.DocumentID != nil},
		EvalID:           uuid.NullUUID{UUID: ptrToUUID(req.Target.EvalID), Valid: req.Target.EvalID != nil},
		EvalItemID:       uuid.NullUUID{UUID: ptrToUUID(req.Target.EvalItemID), Valid: req.Target.EvalItemID != nil},
		AttemptID:        uuid.NullUUID{UUID: ptrToUUID(req.Target.AttemptID), Valid: req.Target.AttemptID != nil},
		Text:             outputText,
		OutputJSON:       outputJSON,
		Model:            modelName,
		Prompt:           promptText,
		PromptRender:     promptText,
		PromptTemplateID: uuid.NullUUID{UUID: promptTmplID, Valid: promptTmplID != uuid.Nil},
		SchemaTemplateID: uuid.NullUUID{UUID: schemaTmplID, Valid: schemaTmplID != uuid.Nil},
		ModelParams:      modelParams,
		Meta:             meta,
		Error:            errorMsg,
	}

	art, err := s.artifactsService.CreateArtifact(ctx, params)
	if err != nil {
		return uuid.Nil, err
	}
	return art.ID, nil
}

func ptrToUUID(p *uuid.UUID) uuid.UUID {
	if p == nil {
		return uuid.Nil
	}
	return *p
}

func modelNameForArtifact(config *ModelConfig) string {
	if config != nil && config.Name != "" {
		return config.Name
	}
	return ""
}

func ptr[T any](v T) *T {
	return &v
}

type artifactMeta struct {
	SystemInstructionID   *uuid.UUID `json:"system_instruction_id,omitempty"`
	SystemInstructionText string     `json:"system_instruction_text,omitempty"`
	ModelConfigID         uuid.UUID  `json:"model_config_id,omitempty"`
}

func buildArtifactMetadata(req GenerateRequest, model *ModelConfig, systemInstruction string) (json.RawMessage, json.RawMessage, error) {
	var modelParams json.RawMessage
	if model != nil {
		serialized, err := json.Marshal(model)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal model params: %w", err)
		}
		modelParams = serialized
	}

	meta := artifactMeta{
		SystemInstructionID:   req.Instructions.SystemInstructionID,
		SystemInstructionText: systemInstruction,
		ModelConfigID:         req.ModelConfigID,
	}

	if meta.SystemInstructionID == nil && meta.SystemInstructionText == "" && meta.ModelConfigID == uuid.Nil {
		return modelParams, nil, nil
	}

	serialized, err := json.Marshal(meta)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal artifact meta: %w", err)
	}

	return modelParams, serialized, nil
}
