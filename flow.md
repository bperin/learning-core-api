                          ┌─────────────────────────┐
                          │      Document Upload     │
                          │   (pdf, text, metadata)  │
                          └──────────────┬──────────┘
                                         │
                                         ▼
                         ┌────────────────────────────────┐
                         │   DocumentReference (DB)       │
                         │ - subject, curriculum, grade   │
                         │ - extracted text               │
                         │ - metadata (admin-assigned)    │
                         └────────────────────────────────┘
                                         │
                                         ▼
        ┌────────────────────────────────────────────────────────────────────┐
        │  STEP 1: INTENT EXTRACTION                                         │
        │────────────────────────────────────────────────────────────────────│
        │  Input: DocumentReference + PromptTemplate(INTENTS) + Schema(INTENTS)  │
        │  AI Output: { intents: [...] }                                      │
        │  Creates Artifact(INTENTS)                                          │
        └───────────────────────┬────────────────────────────────────────────┘
                                │
                                ▼
        ┌────────────────────────────────────────────────────────────────────┐
        │  STEP 2: PLANNING (OPTIONAL)                                       │
        │────────────────────────────────────────────────────────────────────│
        │  Input: DocumentReference + Intents + PromptTemplate(PLAN)         │
        │         + Schema(PLAN)                                             │
        │  AI Output: Plan (sections, item counts, topic coverage)           │
        │  Creates Artifact(PLAN)                                            │
        └───────────────────────┬────────────────────────────────────────────┘
                                │
                                ▼
        ┌────────────────────────────────────────────────────────────────────┐
        │  STEP 3: EVAL GENERATION                                          │
        │────────────────────────────────────────────────────────────────────│
        │  Input: DocumentReference + (Intents) + (Plan)                     │
        │         + PromptTemplate(EVAL) + Schema(EVAL)                      │
        │                                                                    │
        │  AI Output Schema:                                                 │
        │   {                                                                │
        │     "title": "string",                                             │
        │     "items": [                                                     │
        │       { "prompt": "...", "options": [...], "correct_index": N, ... } │
        │     ]                                                              │
        │   }                                                                │
        │                                                                    │
        │  Creates:                                                          │
        │   • Artifact(EVAL)                                                 │
        └───────────────────────┬────────────────────────────────────────────┘
                                │
                                ▼
        ┌────────────────────────────────────────────────────────────────────┐
        │  DOMAIN OBJECT CREATION (IMMUTABLE TRUTH)                          │
        │────────────────────────────────────────────────────────────────────│
        │  Converts raw AI JSON →                                            │
        │   • Eval (immutable)                                               │
        │   • EvalItems[] (immutable)                                        │
        │                                                                    │
        │  These become the student-facing truth.                            │
        └───────────────────────┬────────────────────────────────────────────┘
                                │
                                ▼
        ┌────────────────────────────────────────────────────────────────────┐
        │  STEP 4: STUDENT RUNTIME                                           │
        │────────────────────────────────────────────────────────────────────│
        │  TestAttempt created → UserAnswers written → scoring performed     │
        │  (Eval + EvalItems never mutate)                                   │
        └───────────────────────┬────────────────────────────────────────────┘
                                │
                                ▼
        ┌────────────────────────────────────────────────────────────────────┐
        │  STEP 5: REAL-TIME HINTING                                         │
        │────────────────────────────────────────────────────────────────────│
        │  Input: EvalItem + PromptTemplate(HINT) + Schema
