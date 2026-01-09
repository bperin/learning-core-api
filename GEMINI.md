# Development Guide

## Prerequisites

Before you begin, ensure you have the following installed on your system:

-   **Go**: Latest version
-   **Make**: Standard build tool
-   **PostgreSQL**: Local or remote instance

## Running the Application

To run the Go API:

```bash
make run
```

Or directly:

```bash
go run cmd/api/main.go
```

## Database Migrations

Migrations are managed using `goose` and located in [`internal/persistance/migrations/`](internal/persistance/migrations/).

### Creating Migrations

-   **Rule:** Never update or rewrite an existing migration file. Always create a new one to apply changes.
-   To create a new migration, add a new `.sql` file in the migrations directory following the naming convention (e.g., `0010_new_feature.sql`).

### Applying Migrations

To migrate the database to the latest version:

```bash
make migrate-up
```

This executes `goose up` using the configuration defined in the [`Makefile`](Makefile).

## Persistence Layer & SQLC

The persistence layer is located in [`internal/persistance/`](internal/persistance/). It consists of:

-   **Migrations:** SQL schema changes.
-   **Queries:** SQL query definitions in [`internal/persistance/queries/`](internal/persistance/queries/).
-   **Store:** Generated Go code in [`internal/persistance/store/`](internal/persistance/store/).

### Updating Queries

1. Edit or add `.sql` files in [`internal/persistance/queries/`](internal/persistance/queries/).
2. Run `sqlc` to regenerate the store:
    ```bash
    make sqlc
    ```
3. The generated code in [`internal/persistance/store/`](internal/persistance/store/) is **immutable and read-only**. Do not manually edit files in this directory.

The [`sqlc.yaml`](sqlc.yaml) file defines the mapping between queries, schema, and generated output paths.
