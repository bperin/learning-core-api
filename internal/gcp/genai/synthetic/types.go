package synthetic

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type DocumentReference struct {
	ID          uuid.UUID
	SubjectID   *uuid.UUID
	SubjectName string
	Curriculum  string
	HumanTitle  string
	TopicTags   []string
	ReviewerID  *uuid.UUID
	RequestedBy uuid.UUID
	RequestedAt time.Time
}

type PromptTemplate struct {
	ID          uuid.UUID
	Key         string
	Version     int32
	Template    string
	Model       string
	ModelParams map[string]any
}

type SchemaTemplate struct {
	ID         uuid.UUID
	SchemaType string
	Version    int32
	SchemaJSON json.RawMessage
}

type Artifact struct {
	Type             string
	Status           string
	EvalID           *uuid.UUID
	EvalItemID       *uuid.UUID
	AttemptID        *uuid.UUID
	ReviewerID       *uuid.UUID
	Text             *string
	OutputJSON       json.RawMessage
	Model            *string
	Prompt           *string
	PromptTemplateID *uuid.UUID
	SchemaTemplateID *uuid.UUID
	ModelParams      map[string]any
	PromptRender     *string
	InputHash        *string
	Meta             map[string]any
	Error            *string
	CreatedAt        time.Time
}

type Intent struct {
	Domain               string               `json:"domain"`
	Subject              string               `json:"subject"`
	IntendedAudience     string               `json:"intended_audience"`
	AssumedPrerequisites []string             `json:"assumed_prerequisites"`
	LearningObjectives   []string             `json:"learning_objectives"`
	KeyConcepts          []string             `json:"key_concepts"`
	DifficultyLevel      string               `json:"difficulty_level"` // "introductory" | "intermediate" | "advanced"
	RecommendedArtifacts RecommendedArtifacts `json:"recommended_artifacts"`
}

type RecommendedArtifacts struct {
	Flashcards              int `json:"flashcards"`
	MultipleChoiceQuestions int `json:"multiple_choice_questions"`
	ShortAnswerQuestions    int `json:"short_answer_questions"`
}

type Plan struct {
	Title string     `json:"title"`
	Steps []PlanStep `json:"steps"`
}

type PlanStep struct {
	Title      string   `json:"title"`
	Objectives []string `json:"objectives"`
}

type EvalPayload struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Items       []EvalItemPayload `json:"items"`
}

type EvalItemPayload struct {
	Prompt       string   `json:"prompt"`
	Options      []string `json:"options"`
	CorrectIndex int32    `json:"correct_index"`
	Hint         string   `json:"hint"`
	Explanation  string   `json:"explanation"`
}
