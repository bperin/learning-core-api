# Memorang Data Flow & Ownership Policy

This document defines **roles**, **ownership**, **immutability rules**, and **data flow** for the Memorang platform.

Roles are:

-   **System**
-   **Admin**
-   **Teacher**
-   **Learner**

---

## 1. Roles

### 1.1 System

**Purpose:** Execute policies and enforce invariants.

System can:

-   Run generation pipelines
-   Call models and produce Artifacts
-   Apply active `PromptTemplate` versions
-   Enforce immutability of domain objects
-   Compute correctness metrics and regressions
-   Persist provenance and audit trails

System **cannot**:

-   Modify published `Eval` / `EvalItem`
-   Override Admin/Teacher governance decisions

---

### 1.2 Admin

**Purpose:** Own policy, generation, and lifecycle.

Admin can:

-   Create and version `PromptTemplate`
-   Define generation constraints and safety rules
-   Trigger content generation (Evals, questions)
-   Run async evals (e.g., GCP evals)
-   Inspect metrics and regression reports
-   Promote or roll back prompt versions
-   Configure review queues and workflows for Teachers

Admin **cannot**:

-   Edit published `Eval` / `EvalItem`
-   Modify `TestAttempt` or `UserAnswer` data

---

### 1.3 Teacher

**Purpose:** Own expert review and correctness validation.

Teacher can:

-   Review `EvalItem` quality
-   Submit `EvalItemReview` verdicts
-   Provide expert explanations and corrections via review channels
-   Flag problematic items for Admin action (e.g., new generation runs)

Teacher **cannot**:

-   Create or version `PromptTemplate`
-   Trigger generation or eval pipelines
-   Edit `Eval` / `EvalItem`
-   Edit Artifacts

---

### 1.4 Learner

**Purpose:** Take assessments and receive feedback.

Learner can:

-   Start tests (`TestAttempt`)
-   Submit answers (`UserAnswer`)
-   Receive scores, explanations, and hints

Learner **cannot**:

-   Generate or edit content
-   View prompts, Artifacts, or eval metrics
-   Change policy, prompts, or correctness signals

---

## 2. Immutable Domain Objects (Learner-Facing Truth)

These define **what Learners see**. Once published, they are immutable.

### 2.1 Eval

-   Complete assessment (test)
-   Lifecycle: draft → published → archived

**Owner:** System (generated under Admin policy)
**Mutable:** ❌ after publish
**Visible to Learner:** ✅

---

### 2.2 EvalItem

-   Single question on an Eval
-   Includes stem, options, correct answer, explanation, hint

**Owner:** System
**Mutable:** ❌ after publish
**Visible to Learner:** ✅
**Reviewed by:** Teacher (primary), Admin (override/escalation)

---

### 2.3 TestAttempt

-   One Learner taking one Eval at a point in time
-   Aggregates answers and scoring

**Owner:** Learner
**Mutable:** ⚠️ only while in progress (not submitted)
**Visible to Learner:** ✅
**Visible to Teacher/Admin:** ✅ (for analytics, not mutation)

---

### 2.4 UserAnswer

-   One answer to one EvalItem inside a TestAttempt
-   Includes correctness and timing metadata

**Owner:** Learner
**Mutable:** ❌ after submission
**Visible to Learner:** ✅
**Visible to Teacher/Admin:** ✅ (for analysis and review)

---

## 3. Mutable Governance & Evidence Objects

These describe **how content was produced and evaluated**.
They never change learner-facing truth.

### 3.1 PromptTemplate

-   Versioned prompt policy for:

    -   `eval_generation`
    -   `hint_generation`
    -   `intent_extraction`
    -   other generation tasks

-   Stored in DB (not hard-coded in services)
-   Exactly one active version per key

**Owner:** Admin
**Mutable:** ✅ (versioned; old versions kept)
**Visible to Learner:** ❌
**Visible to Teacher:** ✅ (read-only, if needed for review context)

---

### 3.2 Artifact

-   Persisted record of all AI outputs and provenance

Includes (non-exhaustive):

-   `INTENTS` (document analysis)
-   `PLAN` (generation plans)
-   `EVAL` / `EVAL_ITEM_RAW` (raw model outputs)
-   `PROMPT_RENDER` (prompt text as sent to model)
-   `QUALITY_METRICS` (GCP eval results, scores)
-   `HINT` (runtime hint outputs)

**Owner:** System / Admin
**Mutable:** ✅ append-only (no destructive edits)
**Visible to Learner:** ❌
**Visible to Teacher:** ✅ (read-only for review and auditing)

---

### 3.3 EvalItemReview

-   Human expert judgment on `EvalItem` quality
-   Provides ground-truth correctness and pedagogical quality

Verdicts (examples):

-   `APPROVED`
-   `REJECTED`
-   `NEEDS_REVISION`
-   `AMBIGUOUS`
-   `OUT_OF_SCOPE`

**Owner:** Teacher (primary) / Admin (escalation or override)
**Mutable:** ✅ (can add new reviews and updated verdicts, with audit trail)
**Visible to Learner:** ❌
**Used by System:** ✅ (as signal for future generation and regression monitoring)

---

## 4. Data Flows

### 4.1 Content Generation (Admin → System)

```text
Admin configures PromptTemplate
  ↓
Source Documents ingested
  ↓
System: Intent Analysis → Artifact(INTENTS)
  ↓
System: (optional) Planning → Artifact(PLAN)
  ↓
System: Eval Generation using PromptTemplate
  ↓
Eval + EvalItems created as immutable domain objects
  ↓
Artifacts(EVAL, EVAL_ITEM_RAW, PROMPT_RENDER) persisted
```

-   System executes; Admin triggers and configures.
-   Teacher does not generate; they review after the fact.
-   Learners never see any Artifacts or PromptTemplates.

---

### 4.2 Correctness & Regression Monitoring (Teacher + Admin + System)

```text
EvalItems (published)
  ├─ Teacher Review → EvalItemReview (Teacher-owned)
  └─ System GCP Evals → Artifact(QUALITY_METRICS)
         ↓
      Admin reviews metrics + patterns
         ↓
      Admin adjusts PromptTemplate or constraints
```

-   `Eval` / `EvalItem` remain immutable.
-   Teacher and System provide signals; Admin updates policy.

---

### 4.3 Policy Evolution (Admin, informed by Teacher/System)

```text
EvalItemReview + QUALITY_METRICS
  ↓
Admin analyzes regressions / failure modes
  ↓
Admin updates or creates new PromptTemplate version
  ↓
Future generation uses new version
  ↓
Old evals remain unchanged; new evals follow new policy
```

-   Teacher feedback is a primary correctness source.
-   System metrics are secondary but scalable.
-   Admin is the only role that changes policy.

---

### 4.4 Learner Runtime Flow

```text
Learner starts Test
  ↓
System creates TestAttempt
  ↓
Learner submits UserAnswers
  ↓
System computes score + feedback
  ↓
Learner sees score, explanation, and hints (if enabled)
```

-   Only `TestAttempt` is mutable while in progress.
-   `UserAnswer` is immutable after submission.
-   Admin/Teacher can see results but not mutate them.

---

### 4.5 Real-Time Hinting

```text
EvalItem (immutable)
 + Active hint PromptTemplate
  ↓
System generates hint → Artifact(HINT)
  ↓
Hint delivered to Learner during TestAttempt
```

Rules:

-   Hints are scoped to a single `EvalItem`.
-   Hints do not create new EvalItems or modify correctness.
-   Hints are not treated as new domain truth; they are guidance only.

---

## 5. Mutability Matrix

| Entity         | Mutable                   | Primary Owner   | Other Roles Involved        |
| -------------- | ------------------------- | --------------- | --------------------------- |
| Eval           | ❌ after publish          | System          | Admin (draft/publish flow)  |
| EvalItem       | ❌ after publish          | System          | Teacher/Admin (review-only) |
| TestAttempt    | ⚠️ while in progress      | Learner         | System                      |
| UserAnswer     | ❌ after submission       | Learner         | System                      |
| PromptTemplate | ✅ (versioned)            | Admin           | System (reads), Teacher(ro) |
| Artifact       | ✅ append-only            | System / Admin  | Teacher(ro), Admin(ro/w)    |
| EvalItemReview | ✅ (new reviews, updates) | Teacher / Admin | System (reads as signal)    |

(ro = read-only)

---

## 6. Generation & Governance Invariants

-   Learners **never generate** content.
-   Teachers **validate** content; they do not generate or mutate domain truth.
-   Admins **govern policy**, not learner data or published truth.
-   System **enforces immutability** and executes generation under Admin policy.
-   Published `Eval` and `EvalItem` are **never mutated**.
-   All model calls produce persisted **Artifacts**.
-   Prompt policies live in the database as `PromptTemplate`, **not hard-coded**.
-   Async evals and regression checks **do not mutate** production domain objects.
-   Policy changes apply only to **future** content; existing Evals remain stable.

---

## 7. Role Responsibility Summary

**System**

-   Executes generation, evals, and hinting
-   Persists Artifacts and provenance
-   Enforces immutability guarantees

**Admin**

-   Owns PromptTemplate and policy
-   Owns generation triggers and lifecycle
-   Owns regression management and rollbacks

**Teacher**

-   Owns correctness reviews via EvalItemReview
-   Provides domain-expert validation and feedback
-   Feeds signals into policy evolution (through Admin/System)

**Learner**

-   Owns TestAttempts and UserAnswers
-   Consumes assessments, explanations, and hints
-   Provides behavioral data, not content or policy

---

## 8. One-Line Mental Model

> **System executes, Admin governs, Teacher validates, Learners learn.**
