# Agent Instructions: Persistence & Migrations

When working with the database layer in this repository, follow these strict guidelines.

## Prerequisites

Ensure your environment has the following installed:

-   **Go** (latest version)
-   **Make**

## Database Migrations

-   **Location:** [`internal/persistance/migrations/`](internal/persistance/migrations/)
-   **Immutability:** Once a migration is committed, **never** modify or rewrite it.
-   **Workflow:** To change the schema, always create a new migration file with the next sequential prefix (e.g., `0010_...sql`).
-   **Execution:** Run migrations using:
    ```bash
    make migrate-up
    ```

## SQL Queries and Code Generation

-   **Queries:** SQL files are located in [`internal/persistance/queries/`](internal/persistance/queries/).
-   **Store:** The generated Go code is in [`internal/persistance/store/`](internal/persistance/store/).
-   **Workflow:**
    1. Update or create SQL queries in the queries directory.
    2. Run the generation tool:
        ```bash
        make sqlc
        ```
    3. Verify the changes in the store directory.
-   **Constraints:** The `store` package is **read-only**. Do not manually edit any file within [`internal/persistance/store/`](internal/persistance/store/). All changes must come from `sqlc` generation.

## Configuration

-   [`sqlc.yaml`](sqlc.yaml) controls the generation process. It defines the paths for schema, queries, and output.
-   [`Makefile`](Makefile) contains the orchestration commands for development.
