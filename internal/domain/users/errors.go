package users

import "errors"

// Domain errors for users
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrPasswordTooShort  = errors.New("password too short")
	ErrPasswordTooWeak   = errors.New("password too weak")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrUserInactive      = errors.New("user account is inactive")
	ErrInvalidRole       = errors.New("invalid user role")
	ErrRoleNotFound      = errors.New("user role not found")
	ErrCannotDeleteSelf  = errors.New("cannot delete your own account")
	ErrAdminRequired     = errors.New("admin privileges required")
)
