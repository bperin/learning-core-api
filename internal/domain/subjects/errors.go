package subjects

import "errors"

// Domain errors for subjects
var (
	ErrSubjectNotFound      = errors.New("subject not found")
	ErrInvalidSubjectName   = errors.New("invalid subject name")
	ErrSubjectNameTooLong   = errors.New("subject name too long")
	ErrInvalidDescription   = errors.New("invalid subject description")
	ErrInvalidUserID        = errors.New("invalid user ID")
	ErrUnauthorized         = errors.New("unauthorized access to subject")
	ErrSubjectInUse         = errors.New("subject is in use and cannot be deleted")
	ErrDuplicateSubjectName = errors.New("subject name already exists for this user")
)
