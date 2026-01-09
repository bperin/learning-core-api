# Agent Instructions: Core Practices, Data Flow, Migrations, Testing

This is the single source of truth for working in this repository. If another doc conflicts, follow this file.

## Prerequisites

-   **Go** (latest)
-   **Make**
-   **PostgreSQL**

## Source of Truth (Schema)

-   The Go migrations and SQLC queries are the authoritative model.
-   Use `internal/persistance/migrations/` + `internal/persistance/queries/` as the contract when shaping Go models and services.

## Project Data Flow (High Level)

1. **HTTP**: `internal/http` routes requests into handlers.
2. **Domain**: handlers call domain services in `internal/domain/...`.
3. **Persistence**: services call repositories (SQLC) in `internal/persistance/store`.
4. **Database**: repositories execute SQL against Postgres.
5. **External services**: specialized flows (e.g., file search) live in `internal/filesearch` and coordinate subject/document repos.

## Go Best Practices (Project Standards)

-   Use `context.Context` on every handler/service/repo boundary.
-   Prefer domain services for business rules; keep repositories as thin persistence layers.
-   Keep structs aligned with DB columns and JSON payloads; avoid implicit defaults.
-   Wrap errors with context and return them; do not log and swallow.
-   Avoid global state; inject dependencies explicitly.
-   Keep time handling in UTC (`time.Time`), and prefer `TIMESTAMPTZ` in SQL.

## Database Migrations

-   **Location:** `internal/persistance/migrations/`
-   **Immutability:** Never edit an existing migration once committed.
-   **Workflow:** Create a new migration with the next sequential prefix (e.g., `0010_...sql`).
-   **Execution:**
    ```bash
    make migrate-up
    ```
-   **Destructive changes:** If a migration resets or drops tables, document it clearly in the migration file header.

## SQL Queries and Code Generation

-   **Warning:** SQLC generates the `store` package; never edit generated files by hand.
-   **Queries:** `internal/persistance/queries/`
-   **Generated store:** `internal/persistance/store/`
-   **SQLC is the generator:** treat everything under `internal/persistance/store/` as generated output.
-   **Workflow:**
    1. Update or create SQL queries in the queries directory.
    2. Run:
        ```bash
        make sqlc
        ```
    3. Verify generated changes in the store directory.
-   **Constraint:** The `store` package is read-only; never edit generated files manually.
-   **Config:** `sqlc.yaml` defines schema + query inputs and the output path.

## Testing Strategy (Contract)

### Core Principles

-   Test behavior, not implementation.
-   Prefer real systems over mocks (real DB, HTTP, migrations, constraints).
-   Immutability is a first-class invariant (published content and attempts are append-only).
-   Black-box tests provide the highest value.

### Layers

| Layer      | Goal              | Speed  | Tools                    |
| ---------- | ----------------- | ------ | ------------------------ |
| Repository | Data correctness  | Fast   | `testing`, real DB       |
| Service    | Business rules    | Medium | `testing`, minimal mocks |
| End-to-End | System guarantees | Slower | `httptest`, real DB      |

### Repository Testing

-   Validate schema correctness, constraints, and immutability.
-   Inserts work, invalid writes fail, deletes cascade correctly.
-   Avoid testing HTTP/prompt/AI logic here.

### End-to-End Testing

-   Exercise the HTTP layer with `httptest`.
-   Verify persisted side effects and emitted artifacts via DB reads.
