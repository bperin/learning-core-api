# testing.md — Go Testing Strategy

This document defines **how we test in Go**, what we test at **each layer**, and which tools/patterns we use.  
It is written to support **immutability, correctness, and regression detection**, especially for AI-assisted systems.

This is not theoretical guidance — it is the **operational testing contract**.

---

## 1. Core Principles

### 1.1 Test behavior, not implementation

-   Tests assert **what the system does**, not how it does it
-   Refactors should not break tests unless behavior changes

### 1.2 Prefer real systems over mocks

-   Real DB
-   Real HTTP server
-   Real migrations
-   Real constraints

Mocks are used **only** where the real thing is impossible or too expensive.

---

### 1.3 Immutability is a first-class invariant

We explicitly test that:

-   published evals cannot be modified
-   eval items are immutable
-   artifacts are append-only
-   attempts and answers are never rewritten

If immutability breaks, tests must fail.

---

### 1.4 Black-box tests are the highest value

Especially for:

-   AI generation
-   prompt changes
-   RAG metadata projection
-   regression detection

We test through:

-   HTTP
-   DB reads
-   emitted artifacts

Not internal functions.

---

## 2. Testing Layers Overview

| Layer      | Goal              | Speed  | Tools                    |
| ---------- | ----------------- | ------ | ------------------------ |
| Repository | Data correctness  | Fast   | `testing`, real DB       |
| Service    | Business rules    | Medium | `testing`, minimal mocks |
| End-to-End | System guarantees | Slower | `httptest`, real DB      |

---

## 3. Repository-Level Testing

### Purpose

Validate **data persistence rules**:

-   schema correctness
-   constraints
-   immutability
-   joins
-   indexes

Repositories must not contain business logic.

---

### What we test

-   inserts work
-   invalid writes fail
-   updates are blocked where expected
-   deletes cascade correctly
-   append-only behavior holds

---

### What we do NOT test

-   HTTP behavior
-   prompt logic
-   AI behavior
-   authorization

---

### Pattern

-   Use a **real Postgres database**
-   Run migrations before tests
-   Clean state per test or per package

---

### Example

```go
func TestEvalItemsAreImmutable(t *testing.T) {
  repo := NewEvalItemRepo(testDB)

  item := createEvalItem(t, repo)

  err := repo.UpdatePrompt(item.ID, "new prompt")
  require.Error(t, err)
}
```
