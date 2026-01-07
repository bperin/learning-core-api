learning-engine/
  cmd/
    api/
      main.go

  internal/
    modules/
      controller.go      // HTTP handlers
      service.go         // domain logic
      repository.go      // DB access
      model.go           // entities + value objects

    documents/
      controller.go
      service.go
      repository.go
      model.go

    generation/
      service.go         // enqueue generation, persist candidates
      repository.go
      agents.go          // generator agents only

    evals/
      service.go         // gating + regression
      repository.go
      judges.go          // judge prompts

    sessions/
      controller.go
      service.go
      repository.go
      model.go

    tutor/
      controller.go
      service.go         // runtime hint/teacher agents
      agents.go

    users/
      controller.go
      service.go
      repository.go
      model.go

    infrastructure/
      db/
        postgres.go
      llm/
        gemini.go
        prompts/
      filesearch/
        filesearch.go
      config/
      logging/

  sql/
    schema.sql
