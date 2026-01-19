# Handler Annotations and Role-Based Routing

This repo uses swagger-style godoc comments plus role-based route registration.

## Godoc annotations for handlers

Use the standard swagger annotations on each handler method:

```go
// GetThing godoc
// @Summary Short action name
// @Description Include role access here (e.g., "Admin-only", "Teacher+Learner").
// @Tags things
// @Accept json
// @Produce json
// @Param id path string true "Thing ID"
// @Success 200 {object} things.Thing
// @Failure 400 {string} string "invalid request"
// @Security OAuth2Auth[read]
// @Router /things/{id} [get]
```

Notes:

- `@Security` should match the `@securityDefinitions` name in `cmd/api/main.go`.
- If a handler is public, omit `@Security` and call that out in `@Description`.
- Roles are documented in `@Description`; enforcement happens in middleware.
- Use scope suffixes for read/write: `OAuth2Auth[read]` or `OAuth2Auth[write]`.

## Role-based route registration

Handlers declare routes by role via the registrar interface in `internal/infra/http.go`:

```go
type RoleRouteRegistrar interface {
	RegisterPublicRoutes(r chi.Router)
	RegisterAdminRoutes(r chi.Router)
	RegisterTeacherRoutes(r chi.Router)
	RegisterLearnerRoutes(r chi.Router)
}
```

Registration is centralized in `internal/infra/http.go`:

- `RegisterPublicRoutes` are mounted without auth.
- `RegisterAdminRoutes`, `RegisterTeacherRoutes`, `RegisterLearnerRoutes` are wrapped with:
    - `JWTMiddleware(secret)`
    - `RequireRole(role)`

Use this pattern in each handler to keep access rules explicit and consistent with
`FLOW_OWNERSHIP.md`.

## JWT role and scope enforcement

Roles are extracted from the JWT `roles` claim in `internal/infra/middleware.go`
and stored in context by `internal/http/authz`.
Use these constants when describing access:

```go
const (
	RoleAdmin   = "admin"
	RoleTeacher = "teacher"
	RoleLearner = "learner"
)
```

Scopes are extracted from `scopes` (array or space-delimited string), or
from `scope` (space-delimited string). If no scope claim is provided,
the middleware infers scopes from roles (`read` for any authenticated user,
`write` for admin/teacher/learner).

Use `RequireScope("read")` and `RequireScope("write")` in handlers to enforce
scope requirements per route.

## Schema Management Handlers

The following handlers manage system configuration and generation policies:

### Prompt Templates & Schema Templates

- **Admin**: Full CRUD access (create, list, get, activate, deactivate)
- **Teacher**: Read-only access to active versions via `GetActiveByGenerationType`
- **Learner**: No access (routes not registered)
- **Security**: `@Security OAuth2Auth[read]` for GET operations, `@Security OAuth2Auth[write]` for POST operations

### Chunking Configs & System Instructions

- **Admin**: Full CRUD access (create, list, get, activate)
- **Teacher**: No access (routes not registered)
- **Learner**: No access (routes not registered)
- **Security**: `@Security OAuth2Auth[read]` for GET operations, `@Security OAuth2Auth[write]` for POST operations

All schema management endpoints enforce immutability: new versions are created rather than edited, and only admin can activate/deactivate versions.
