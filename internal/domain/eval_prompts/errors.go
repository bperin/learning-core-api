package eval_prompts

import "fmt"

var (
	ErrInvalidEvalType  = fmt.Errorf("eval type is required")
	ErrEmptyPromptText  = fmt.Errorf("prompt text is required")
	ErrPromptNotFound   = fmt.Errorf("eval prompt not found")
	ErrInvalidVersion   = fmt.Errorf("invalid prompt version")
	ErrNoActivePrompt   = fmt.Errorf("no active prompt found for eval type")
)
