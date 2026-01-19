package seeds

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"

	"learning-core-api/internal/persistance/store"
	"learning-core-api/internal/utils"
)

const (
	systemInstructionsSeed = "system_instructions.txt"
	taxonomyPromptSeed     = "taxonomy_prompt.txt"
	taxonomySchemaSeed     = "taxonomy_schema.json"
	questionsPromptSeed    = "questions_prompt.txt"
	questionsSchemaSeed    = "questions_schema.json"
	chunkingConfigSeedFile = "chunking_config.json"

	systemSeedEmail    = "admin@test.local"
	systemSeedPassword = "seed_placeholder_password"
)

type promptTemplateSeed struct {
	GenerationType string
	Title          string
	Description    *string
	Template       string
	Metadata       json.RawMessage
	CreatedBy      *string
}

type schemaTemplateSeed struct {
	GenerationType string
	SchemaJSON     json.RawMessage
	IsActive       *bool
}

type chunkingConfigSeedPayload struct {
	ChunkingConfig struct {
		WhiteSpaceConfig struct {
			MaxTokensPerChunk int32 `json:"max_tokens_per_chunk"`
			MaxOverlapTokens  int32 `json:"max_overlap_tokens"`
		} `json:"white_space_config"`
	} `json:"chunking_config"`
}

func Run(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is required")
	}
	return RunWithQueries(ctx, store.New(db))
}

func RunWithQueries(ctx context.Context, queries *store.Queries) error {
	if queries == nil {
		return fmt.Errorf("queries are required")
	}

	systemUserID, err := ensureSystemUser(ctx, queries)
	if err != nil {
		return fmt.Errorf("failed to ensure system user: %w", err)
	}

	if err := seedSystemInstructions(ctx, queries, systemUserID); err != nil {
		return fmt.Errorf("failed to seed system instructions: %w", err)
	}

	if err := seedModelConfig(ctx, queries, systemUserID); err != nil {
		return fmt.Errorf("failed to seed model configs: %w", err)
	}

	if err := seedChunkingConfig(ctx, queries, systemUserID); err != nil {
		return fmt.Errorf("failed to seed chunking configs: %w", err)
	}

	if err := seedPromptTemplates(ctx, queries); err != nil {
		return fmt.Errorf("failed to seed prompt templates: %w", err)
	}

	if err := seedSchemaTemplates(ctx, queries, systemUserID); err != nil {
		return fmt.Errorf("failed to seed schema templates: %w", err)
	}

	if err := seedSubjects(ctx, queries); err != nil {
		return fmt.Errorf("failed to seed subjects: %w", err)
	}

	return nil
}

func seedModelConfig(ctx context.Context, queries *store.Queries, createdBy uuid.UUID) error {
	_, err := queries.GetActiveModelConfig(ctx)
	if err == nil {
		log.Printf("active model config already exists")
		return nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	_, err = queries.CreateModelConfig(ctx, store.CreateModelConfigParams{
		ModelName:   "gemini-3-flash-preview",
		Temperature: sql.NullFloat64{Float64: 0.9, Valid: true},
		MaxTokens:   sql.NullInt32{Int32: 8192, Valid: true},
		TopP:        sql.NullFloat64{Float64: 0.5, Valid: true},
		TopK:        sql.NullFloat64{Float64: 20, Valid: true},
		MimeType:    sql.NullString{String: "application/json", Valid: true},
		IsActive:    true,
		CreatedBy:   createdBy,
	})
	if err != nil {
		return err
	}
	log.Printf("seeded default model config")
	return nil
}

func seedChunkingConfig(ctx context.Context, queries *store.Queries, createdBy uuid.UUID) error {
	_, err := queries.GetActiveChunkingConfig(ctx)
	if err == nil {
		log.Printf("active chunking config already exists")
		return nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	path, err := seedPath(chunkingConfigSeedFile)
	if err != nil {
		return err
	}
	raw, ok, err := readSeedJSON(path)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("chunking config seed not found: %s", path)
	}
	if len(raw) == 0 {
		return fmt.Errorf("chunking config seed is empty: %s", path)
	}

	var seed chunkingConfigSeedPayload
	if err := json.Unmarshal(raw, &seed); err != nil {
		return fmt.Errorf("failed to parse chunking config seed: %w", err)
	}

	maxTokens := seed.ChunkingConfig.WhiteSpaceConfig.MaxTokensPerChunk
	maxOverlap := seed.ChunkingConfig.WhiteSpaceConfig.MaxOverlapTokens
	if maxTokens <= 0 {
		return fmt.Errorf("chunking config max_tokens_per_chunk must be > 0")
	}
	if maxOverlap < 0 {
		return fmt.Errorf("chunking config max_overlap_tokens must be >= 0")
	}

	isActive := true
	_, err = queries.CreateChunkingConfig(ctx, store.CreateChunkingConfigParams{
		ChunkSize:    maxTokens,
		ChunkOverlap: maxOverlap,
		IsActive:     isActive,
		CreatedBy:    createdBy,
	})
	if err != nil {
		return err
	}
	log.Printf("seeded default chunking config")
	return nil
}

func ensureSystemUser(ctx context.Context, queries *store.Queries) (uuid.UUID, error) {
	user, err := queries.GetUserByEmail(ctx, systemSeedEmail)
	if err == nil {
		log.Printf("system user already exists: email=%s", systemSeedEmail)
		return user.ID, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, err
	}

	userID := uuid.New()
	created, err := queries.CreateUser(ctx, store.CreateUserParams{
		ID:        userID,
		Email:     systemSeedEmail,
		Password:  systemSeedPassword,
		IsAdmin:   true,
		IsLearner: false,
		IsTeacher: false,
	})
	if err != nil {
		return uuid.Nil, err
	}

	log.Printf("seeded system user: email=%s", systemSeedEmail)
	return created.ID, nil
}

func seedSystemInstructions(ctx context.Context, queries *store.Queries, createdBy uuid.UUID) error {
	path, err := seedPath(systemInstructionsSeed)
	if err != nil {
		return err
	}
	text, ok, err := readSeedText(path)
	if err != nil {
		return err
	}
	if !ok {
		log.Printf("no system instructions seed found in %s", path)
		return nil
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return fmt.Errorf("system instructions seed is empty: %s", path)
	}

	active, err := queries.GetActiveSystemInstruction(ctx)
	if err == nil && strings.TrimSpace(active.Text) == text {
		log.Printf("system instructions already active")
		return nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	isActive := true
	_, err = queries.CreateSystemInstruction(ctx, store.CreateSystemInstructionParams{
		Text:      text,
		IsActive:  isActive,
		CreatedBy: createdBy,
	})
	if err != nil {
		return err
	}
	log.Printf("seeded system instructions from %s", path)
	return nil
}

func seedPromptTemplates(ctx context.Context, queries *store.Queries) error {
	type promptSeedDefinition struct {
		filename       string
		generationType utils.GenerationType
		title          string
		description    string
	}

	seeds := []promptSeedDefinition{
		{
			filename:       taxonomyPromptSeed,
			generationType: utils.GenerationTypeClassification,
			title:          "Taxonomy Classification Prompt",
			description:    "Seed prompt template for taxonomy classification",
		},
		{
			filename:       questionsPromptSeed,
			generationType: utils.GenerationTypeQuestions,
			title:          "Question Generation Prompt",
			description:    "Seed prompt template for question generation",
		},
	}

	for _, def := range seeds {
		path, err := seedPath(def.filename)
		if err != nil {
			return err
		}
		promptText, ok, err := readSeedText(path)
		if err != nil {
			return err
		}
		if !ok {
			log.Printf("no prompt template seeds found in %s", path)
			continue
		}

		promptText = strings.TrimSpace(promptText)
		if promptText == "" {
			return fmt.Errorf("prompt template seed is empty: %s", path)
		}

		seed := promptTemplateSeed{
			GenerationType: def.generationType.String(),
			Title:          def.title,
			Description:    stringPtr(def.description),
			Template:       promptText,
			CreatedBy:      stringPtr(systemSeedEmail),
		}

		existing, err := queries.GetPromptTemplatesByGenerationType(ctx, def.generationType.DB())
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			log.Printf("prompt templates already exist: generation_type=%s", seed.GenerationType)
			continue
		}

		normalizedSeedMeta, err := normalizeJSON(seed.Metadata)
		if err != nil {
			return fmt.Errorf("invalid prompt metadata: %w", err)
		}

		for _, tmpl := range existing {
			if tmpl.Title != seed.Title {
				continue
			}
			if tmpl.Template != seed.Template {
				continue
			}
			if tmpl.Description.Valid && seed.Description == nil {
				continue
			}
			if !tmpl.Description.Valid && seed.Description != nil {
				continue
			}
			if tmpl.Description.Valid && seed.Description != nil && tmpl.Description.String != *seed.Description {
				continue
			}

			normalizedExisting, err := normalizeJSON(tmpl.Metadata.RawMessage)
			if err != nil {
				return err
			}
			if jsonEqual(normalizedExisting, normalizedSeedMeta) {
				log.Printf("prompt template already exists: generation_type=%s", seed.GenerationType)
				continue
			}
		}

		_, err = queries.CreateNewVersion(ctx, store.CreateNewVersionParams{
			GenerationType: def.generationType.DB(),
			IsActive:       true,
			Title:          seed.Title,
			Description:    sql.NullString{String: stringValue(seed.Description), Valid: seed.Description != nil},
			Template:       seed.Template,
			Metadata:       toNullRawMessage(seed.Metadata),
			CreatedBy:      sql.NullString{String: stringValue(seed.CreatedBy), Valid: seed.CreatedBy != nil},
		})
		if err != nil {
			return err
		}
		log.Printf("seeded prompt template: generation_type=%s", seed.GenerationType)
	}

	return nil
}

func seedSchemaTemplates(ctx context.Context, queries *store.Queries, createdBy uuid.UUID) error {
	type schemaSeedDefinition struct {
		filename       string
		generationType utils.GenerationType
	}

	seeds := []schemaSeedDefinition{
		{
			filename:       taxonomySchemaSeed,
			generationType: utils.GenerationTypeClassification,
		},
		{
			filename:       questionsSchemaSeed,
			generationType: utils.GenerationTypeQuestions,
		},
	}

	for _, def := range seeds {
		path, err := seedPath(def.filename)
		if err != nil {
			return err
		}
		schemaJSON, ok, err := readSeedJSON(path)
		if err != nil {
			return err
		}
		if !ok {
			log.Printf("no schema template seeds found in %s", path)
			continue
		}
		if len(schemaJSON) == 0 {
			return fmt.Errorf("schema template seed is empty: %s", path)
		}

		seed := schemaTemplateSeed{
			GenerationType: def.generationType.String(),
			SchemaJSON:     schemaJSON,
			IsActive:       boolPtr(true),
		}

		existing, err := queries.ListSchemaTemplatesByGenerationType(ctx, def.generationType.DB())
		if err != nil {
			return err
		}
		if len(existing) > 0 {
			log.Printf("schema templates already exist: generation_type=%s", seed.GenerationType)
			continue
		}

		normalizedSeed, err := normalizeJSON(seed.SchemaJSON)
		if err != nil {
			return fmt.Errorf("invalid schema_json: %w", err)
		}

		for _, tmpl := range existing {
			normalizedExisting, err := normalizeJSON(tmpl.SchemaJson)
			if err != nil {
				return err
			}
			if jsonEqual(normalizedExisting, normalizedSeed) {
				log.Printf("schema template already exists: generation_type=%s", seed.GenerationType)
				continue
			}
		}

		_, err = queries.CreateSchemaTemplate(ctx, store.CreateSchemaTemplateParams{
			GenerationType: def.generationType.DB(),
			SchemaJson:     seed.SchemaJSON,
			IsActive:       boolValue(seed.IsActive),
			CreatedBy:      createdBy,
			LockedAt:       sql.NullTime{},
		})
		if err != nil {
			return err
		}
		log.Printf("seeded schema template: generation_type=%s", seed.GenerationType)
	}

	return nil
}

func boolPtr(value bool) *bool {
	return &value
}

func boolValue(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
}

func normalizeJSON(raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var payload interface{}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	normalized, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return normalized, nil
}

func jsonEqual(a, b json.RawMessage) bool {
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	return string(a) == string(b)
}

func stringPtr(value string) *string {
	return &value
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func toNullRawMessage(raw json.RawMessage) pqtype.NullRawMessage {
	if len(raw) == 0 {
		return pqtype.NullRawMessage{}
	}
	return pqtype.NullRawMessage{RawMessage: raw, Valid: true}
}

func seedPath(filename string) (string, error) {
	dir, err := seedDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, filename), nil
}

func seedDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to determine seeds path")
	}
	return filepath.Dir(file), nil
}

func readSeedText(path string) (string, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, err
	}
	return string(data), true, nil
}

func readSeedJSON(path string) (json.RawMessage, bool, error) {
	text, ok, err := readSeedText(path)
	if err != nil || !ok {
		return nil, ok, err
	}
	return json.RawMessage(text), true, nil
}

func seedSubjects(ctx context.Context, queries *store.Queries) error {
	// Check if subjects already exist
	existingSubjects, err := queries.GetAllSubjects(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if len(existingSubjects) > 0 {
		log.Printf("subjects already seeded (%d subjects found), skipping seed", len(existingSubjects))
		return nil
	}

	// Load subjects from JSON file in seeds directory
	type seedSubject struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}

	path, err := seedPath("subjects.json")
	if err != nil {
		return err
	}

	raw, ok, err := readSeedJSON(path)
	if err != nil {
		return err
	}
	if !ok {
		log.Printf("no subjects seed found in %s", path)
		return nil
	}
	if len(raw) == 0 {
		return fmt.Errorf("subjects seed is empty: %s", path)
	}

	var subjects []seedSubject
	if err := json.Unmarshal(raw, &subjects); err != nil {
		return fmt.Errorf("failed to parse subjects seed: %w", err)
	}

	if len(subjects) == 0 {
		log.Println("no subjects found in seed file, skipping seed")
		return nil
	}

	// Insert subjects into database
	for _, subject := range subjects {
		_, err := queries.CreateSubject(ctx, store.CreateSubjectParams{
			ID:        uuid.New(),
			Name:      subject.Name,
			Url:       fmt.Sprintf("https://open.umn.edu/opentextbooks/subjects/%s", subject.Slug),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})
		if err != nil {
			log.Printf("Warning: failed to create subject %s: %v", subject.Name, err)
			continue
		}
	}

	log.Printf("seeded %d subjects", len(subjects))
	return nil
}
