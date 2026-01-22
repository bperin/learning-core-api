# learning-core-api

Core API for the learning system. This repo defines the authoritative data flow, ownership, and evaluation contracts for generation, review, and learner runtime.

## Where the system is now

- Generation supports three classes: classification (taxonomy extraction), question generation (questions + expected answers), and section topic extraction (section topics with difficulty scores).
- Evaluation is scoped to groundedness and answerability checks against source documents.
- Immutability rules and ownership boundaries are enforced via domain and persistence layers.
- Groundedness evals against GCP are implemented or near implementation.

## What the system aims to be (evals)

The eval surface area is deliberately narrow and normative:

Eval tasks judge outputs against context for factual support and validity. They do not optimize prompts, compare generations, assume gold datasets, or assume determinism.

### Classification evals (taxonomy extraction)

Classification answers: "Does this document support assigning it to these taxonomy nodes?"

Checkable now:
- Concept groundedness (hard): Is the concept explicitly supported by the document text?
- Hallucinated nodes (hard): Did the model invent a concept not present in the document?
- Hierarchy correctness (soft): Does the document imply X is a subtopic of Y? Advisory only.

Not checkable yet:
- Best taxonomy, depth selection, or preferred ontology alignment.

### Question generation evals (questions + expected answers)

Each generated question asserts:
- The question is answerable from the documents.
- The expected answer is correct and supported.

Checkable now:
- Expected-answer groundedness (hard): fully supported by document context.
- Question answerability (hard): answerable using only the provided context.
- Question-answer alignment (soft): useful but warning-only initially.

Not checkable yet:
- Pedagogical quality, difficulty, or coverage completeness.

### Unifying rule

We evaluate whether generated artifacts (taxonomy nodes, questions, expected answers) are factually supported and answerable using the provided documents. We do not evaluate prompt quality, pedagogical optimality, or model performance against static datasets.

## Graph RAG for structured documents

Structured documents (manuals, textbooks, schematics) benefit from graph-aware retrieval because their meaning is embedded in headings, sections, tables, and layout order. This project supports a Graph RAG pipeline where Document AI extracts structure, the system persists a document graph, and generation queries expand context across neighboring nodes.

Pipeline overview:
- Ingest document into GCS and run Document AI OCR/layout processing.
- Convert layout into graph nodes (document/page/paragraph) and edges (contains/next).
- Store graph in Postgres and use `graph_rag` to expand retrieval context.
- Feed the graph context into generation prompts for structure-aware grounding.

Why this helps:
- Preserves section boundaries, ordering, and hierarchy across large documents.
- Enables targeted retrieval of related content instead of flat chunking.
- Provides a bridge to structured downstream artifacts (e.g., CAD outputs).

Long-term: graph-backed RAG enables pipelines that transform structured source documents into raw CAD assets. The graph preserves the structural context required to map diagrams, tables, and references into machine-usable formats during later processing stages.

## Roadmap shape (fan-out and human gating)

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

## Mutability matrix (testing anchor)

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

Testing should enforce immutability guarantees and the eval gates above.

## Source of truth

Database migrations and SQLC queries are the contract:
- `internal/persistance/migrations/`
- `internal/persistance/queries/`

Do not edit generated files in `internal/persistance/store/` by hand.

## Key docs

- `flow.md`: consolidated flow, ownership, and evaluation contract
- `AGENTS.md`: repository rules and testing strategy

## Common commands

- `make sqlc`
- `make migrate-up`
- `make test`
