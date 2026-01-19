# System Flow, Ownership, and Evaluation Contract

This document consolidates system flow, ownership, mutability, and evaluation scope into a single contract.

## 1. What the notebook establishes (normative, not implementation)

The notebook establishes a pattern, not a product feature:

- Eval tasks judge outputs against context for factual support and validity.
- Eval tasks do not optimize prompts.
- Eval tasks do not compare generations to each other.
- Eval tasks do not assume gold datasets.
- Eval tasks do not assume determinism.

This is a strong fit for this system, but only for certain kinds of checks.

## 1.1 Where the system is now

- Generation supports classification (taxonomy extraction), question generation, and section topic extraction with difficulty scoring.
- Evaluation is scoped to groundedness and answerability checks against source documents.
- Groundedness evals against GCP are implemented or near implementation.

## 1.2 Roadmap shape (fan-out and human gating)

Near-term focus (implemented or close):
- Groundedness and answerability evals run against source documents.
- Hard gates on supportedness; soft warnings on alignment and hierarchy.

Nice-to-haves (fan-out pattern):
- Expand eval suites to cover more artifact types and domains.
- Add aggregation for eval suites and trend monitoring.
- Introduce prompt optimization and regression analysis as separate workflows.

Human gating and removal over time:
- Today: human review is required for ambiguous or soft-fail cases.
- Next: use eval warnings to queue targeted reviews rather than full manual review.
- Later: automated publish gates for hard checks; human review only for exceptions.

## 2. Generation classes and eval contracts

Your system has three generation classes with different eval contracts.

### 2.1 Classification (taxonomy extraction)

Classification answers:
"Does this document support assigning it to these taxonomy nodes?"

Each generated taxonomy node is an implicit claim:
- "This concept exists in the document."
- "This concept belongs under this parent."

What you can evaluate now:
- Concept groundedness (hard): Is the concept explicitly supported by the document text?
- Hallucinated nodes (hard): Did the model invent a concept not present in the document?
- Hierarchy correctness (soft): Does the document imply X is a subtopic of Y? Use as advisory only.

What you should not evaluate yet:
- "Is this the best taxonomy?"
- "Did it pick the right depth?"
- "Did it match our preferred ontology?"

Classification eval summary:
- Concept groundedness: Yes (hard gate)
- Hallucinated nodes: Yes (hard gate)
- Hierarchy correctness: Soft (warning only)
- Ontology quality: No (human-only)

### 2.2 Question generation (questions + expected answers)

Each generated question asserts two claims:
- The question is answerable from the documents.
- The expected answer is correct and supported.

What you are not claiming:
- Pedagogical optimality
- Ideal difficulty
- Completeness of coverage

What you can evaluate now:
- Expected-answer groundedness (hard): Is the expected answer fully supported by the documents?
- Question answerability (hard): Can this question be answered using only the provided context?
- Question-answer alignment (soft): Does the expected answer actually answer the question?

What you should not evaluate:
- "Is this a good question?"
- "Is this an interesting question?"
- "Is the difficulty appropriate?"
- "Did we cover all topics?"

Question eval summary:
- Expected-answer groundedness: Yes (hard gate)
- Question answerability: Yes (hard gate)
- QA alignment: Soft (warning only)
- Pedagogical quality: No (human-only)

## 3. Unifying eval principle

Eval tasks only check factual and structural validity of generated artifacts against source documents.
They do not judge creativity, pedagogy, or compare generations.

## 4. Data flow (high level)

```text
Admin configures PromptTemplate and SchemaTemplate
  -> Documents ingested and stored with metadata
  -> System runs generation
     - Classification: taxonomy nodes
     - Questions: questions + expected answers
  -> Artifacts persisted (prompts, outputs, metrics)
  -> Domain objects created as immutable truth (Eval, EvalItem)
  -> Learner runtime: TestAttempt and UserAnswer
  -> Optional hinting uses active prompt + schema, produces Artifact(HINT)
```

## 5. Roles and ownership

Roles:
- System
- Admin
- Teacher
- Learner

### 5.1 System

Purpose: execute policies and enforce invariants.

System can:
- Run generation pipelines
- Apply active prompt and schema templates
- Persist artifacts and provenance
- Enforce immutability

System cannot:
- Modify published Eval or EvalItem

### 5.2 Admin

Purpose: own policy, generation, and lifecycle.

Admin can:
- Create and version PromptTemplate and SchemaTemplate
- Trigger generation and eval runs
- Inspect metrics and regression reports
- Promote or roll back prompt versions

Admin cannot:
- Edit published Eval or EvalItem
- Modify TestAttempt or UserAnswer

### 5.3 Teacher

Purpose: expert review and correctness validation.

Teacher can:
- Review EvalItems
- Submit EvalItemReview verdicts
- Flag issues for Admin action

Teacher cannot:
- Create or version PromptTemplate
- Trigger generation pipelines
- Edit Eval or EvalItem

### 5.4 Learner

Purpose: take assessments and receive feedback.

Learner can:
- Start tests (TestAttempt)
- Submit answers (UserAnswer)
- Receive scores and hints

Learner cannot:
- Generate or edit content
- View prompts, artifacts, or eval metrics

## 6. Domain objects and mutability

### 6.1 Immutable domain objects (learner-facing truth)

- Eval: immutable after publish
- EvalItem: immutable after publish
- UserAnswer: immutable after submission

### 6.2 Mutable governance and evidence objects

- PromptTemplate: versioned and mutable
- SchemaTemplate: versioned and mutable
- Artifact: append-only
- EvalItemReview: additive review history

## 7. Mutability matrix (testing anchor)

| Entity         | Mutable              | Primary Owner   | Notes                         |
| -------------- | -------------------- | --------------- | ----------------------------- |
| Eval           | No (after publish)   | System          | Immutable learner-facing data |
| EvalItem       | No (after publish)   | System          | Immutable learner-facing data |
| TestAttempt    | Soft (in progress)   | Learner         | Mutable until completion      |
| UserAnswer     | No (after submission)| Learner         | Append-only                   |
| PromptTemplate | Yes (versioned)      | Admin           | Policy changes only           |
| SchemaTemplate | Yes (versioned)      | Admin           | Output contract only          |
| Artifact       | Yes (append-only)    | System / Admin  | Provenance and metrics        |
| EvalItemReview | Yes (additive)       | Teacher / Admin | New reviews supersede old     |

Testing should anchor on immutability: immutable rows never mutate; append-only rows only grow; mutable rows only change in allowed windows (TestAttempt).

## 8. Generation and governance invariants

- Learners never generate content.
- Teachers validate content; they do not generate or mutate domain truth.
- Admins govern policy, not learner data or published truth.
- System enforces immutability and executes generation under Admin policy.
- Published Eval and EvalItem are never mutated.
- All model calls produce persisted Artifacts.
- Prompt and schema policies live in the database, not code.

## 9. One-line mental model

System executes, Admin governs, Teacher validates, Learners learn.
