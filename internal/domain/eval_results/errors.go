package eval_results

import "fmt"

var (
	ErrInvalidEvalItemID   = fmt.Errorf("eval item id is required")
	ErrInvalidEvalType     = fmt.Errorf("eval type is required")
	ErrInvalidEvalPromptID = fmt.Errorf("eval prompt id is required")
	ErrInvalidVerdict      = fmt.Errorf("verdict must be PASS, FAIL, or WARN")
	ErrResultNotFound      = fmt.Errorf("eval result not found")
)
