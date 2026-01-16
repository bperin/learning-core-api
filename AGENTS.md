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

# ownership.md — Ownership, Mutability, and Responsibility Model

This document defines **who owns what**, **what is mutable vs immutable**, and **who is allowed to create or change data** across the stack.

It reflects the finalized role model:

-   **Learner** (default)
-   **Teacher**
-   **Admin**

This is a **governance and correctness contract**, not UI guidance.

---

## 1. Roles

### 1.1 Learner (default role)

**Purpose:** Consume learning material and demonstrate knowledge.

Learners can:

-   Take published evaluations
-   Submit answers
-   Receive scores and feedback
-   Request real-time hints on questions

Learners cannot:

-   Create or modify content
-   Generate questions or tests
-   View prompts, artifacts, or eval metrics
-   Influence generation policy

---

### 1.2 Teacher

**Purpose:** Curate and validate educational content.

Teachers can:

-   Upload source documents
-   Assign subject and curriculum metadata
-   Review generated questions
-   Approve or reject questions
-   Trigger re-generation workflows (indirectly)

Teachers cannot:

-   Modify prompt templates
-   Change generation policy
-   Run system-wide evals
-   Override immutability rules

Teachers act as **human correctness validators**, not policy owners.

---

### 1.3 Admin

**Purpose:** Own system policy, generation logic, and correctness monitoring.

Admins can:

-   Create, version, and activate prompt templates
-   Define and update generation constraints (policy)
-   Trigger eval generation jobs
-   Run async correctness and regression evals (GCP)
-   Monitor metrics and drift
-   Promote or roll back prompt versions

Admins do not:

-   Take tests as learners
-   Edit published content inline
-   Mutate immutable domain objects

Admins govern **how content is produced**, not individual outcomes.

---

## 2. Ownership by Layer

### 2.1 Database (Source of Truth)

The database is the **authoritative state** of the system.

-   Stores immutable domain objects
-   Stores append-only artifacts
-   Stores policy as data (not code)

Ownership:

-   Learners: their own attempts and answers
-   Teachers: document metadata and reviews
-   Admins: prompts, policies, and eval generation artifacts

---

### 2.2 Application / Domain Layer (Go)

The domain layer:

-   Enforces immutability
-   Enforces role-based permissions
-   Prevents invalid state transitions

Ownership:

-   System-owned logic
-   No role may bypass domain invariants

---

### 2.3 AI / RAG / External Systems

External systems are **execution engines**, never sources of truth.

-   AI generates candidate content
-   RAG retrieves context based on metadata
-   GCP evals produce metrics

Ownership:

-   Outputs are always persisted back as artifacts
-   External systems never own state

---

## 3. Domain Objects — Ownership & Mutability

### Eval (Test)

-   Represents a published assessment
-   Contains many questions (EvalItems)

Owner: System (generated under admin policy)  
Mutable: ❌ (after publish)  
Who can create: Admin  
Who can view: Learner, Teacher, Admin  
Who can change: Nobody (after publish)

---

### EvalItem (Question)

-   A single question in an Eval

Owner: System  
Mutable: ❌  
Who can create: Admin (via generation)  
Who can review: Teacher, Admin  
Who can answer: Learner  
Who can change: Nobody

---

### TestAttempt

-   A learner taking an Eval at a point in time

Owner: Learner  
Mutable: ⚠️ (only while in progress)  
Who can create: Learner  
Who can complete: Learner  
Who can view: Learner, Admin  
Who can change after completion: Nobody

---

### UserAnswer

-   A learner’s answer to a question

Owner: Learner  
Mutable: ❌  
Who can create: Learner  
Who can view: Learner, Admin  
Who can change: Nobody

---

## 4. Governance & Learning Objects

### Document

-   Source material for generation

Owner: Teacher / Admin  
Mutable: ⚠️ (metadata only)  
Who can upload: Teacher, Admin  
Who can edit metadata: Teacher, Admin  
Who can delete: Admin only

Document content is immutable after ingest.

---

### Subject / Curriculum

-   Human-defined classification labels

Owner: Admin  
Mutable: ✅  
Who can assign: Teacher, Admin  
Who can redefine globally: Admin

These are **classification**, not ontology.

---

### PromptTemplate

-   Versioned generation policy stored in DB

Owner: Admin  
Mutable: ✅  
Who can create/update: Admin  
Who can activate/deprecate: Admin  
Who can view: Admin only

Prompts never live in code.

---

### Artifact

-   Evidence, provenance, metrics, AI outputs

Owner: System  
Mutable: ✅ (append-only)  
Who can create: System, Admin  
Who can view: Admin  
Who can delete: Nobody (except via retention policy)

Artifacts explain **why**, never **what**.

---

### EvalItemReview

-   Human judgment on question correctness

Owner: Teacher / Admin  
Mutable: ✅  
Who can create: Teacher, Admin  
Who can view: Admin  
Who can change: Nobody (new reviews supersede old ones)

This is the **gold correctness signal**.

---

## 5. Mutability Matrix

| Entity         | Mutable | Owner         | Notes                    |
| -------------- | ------- | ------------- | ------------------------ |
| Eval           | ❌      | System        | Immutable after publish  |
| EvalItem       | ❌      | System        | Never changes            |
| TestAttempt    | ⚠️      | Learner       | Mutable until completion |
| UserAnswer     | ❌      | Learner       | Append-only              |
| Document       | ⚠️      | Teacher/Admin | Metadata only            |
| Subject        | ✅      | Admin         | Low-cardinality          |
| Curriculum     | ✅      | Admin         | Low-cardinality          |
| PromptTemplate | ✅      | Admin         | Versioned policy         |
| Artifact       | ✅      | System        | Append-only              |
| EvalItemReview | ✅      | Teacher/Admin | Additive signal          |

---

## 6. Generation Rules (Hard Invariants)

-   Learners never generate content
-   Teachers never change policy
-   Admins never mutate published content
-   Prompts must exist in the DB
-   All AI outputs must be persisted as artifacts
-   External systems never own state
-   Corrections happen via regeneration, not mutation

---

## 7. Responsibility Summary

**Learners**

-   Provide answer data
-   Consume assessments

**Teachers**

-   Curate source material
-   Validate question quality

**Admins**

-   Define generation policy
-   Monitor correctness and regressions
-   Control system evolution

**System**

-   Enforces invariants
-   Persists truth
-   Measures quality
-   Protects learners from experimentation

---

## 8. One-Line Mental Model

> Learners consume truth, teachers validate truth, admins govern policy, and the system enforces immutability.

This document defines the ownership and mutability guarantees of the platform.
