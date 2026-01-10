package evals

import "errors"

// Domain errors for evals
var (
	ErrEvalNotFound            = errors.New("evaluation not found")
	ErrInvalidTitle            = errors.New("invalid evaluation title")
	ErrTitleTooLong            = errors.New("evaluation title too long")
	ErrInvalidDescription      = errors.New("invalid evaluation description")
	ErrInvalidStatus           = errors.New("invalid evaluation status")
	ErrInvalidDifficulty       = errors.New("invalid difficulty level")
	ErrInvalidInstructions     = errors.New("invalid instructions")
	ErrInvalidUserID           = errors.New("invalid user ID")
	ErrUnauthorized            = errors.New("unauthorized access to evaluation")
	ErrCannotModifyPublished   = errors.New("cannot modify published evaluation")
	ErrCannotModifyArchived    = errors.New("cannot modify archived evaluation")
	ErrCannotPublishDraft      = errors.New("cannot publish evaluation in current state")
	ErrCannotArchive           = errors.New("cannot archive evaluation in current state")
	ErrCannotDeletePublished   = errors.New("cannot delete published evaluation")
	ErrEvalHasItems            = errors.New("evaluation has items and cannot be deleted")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
)
