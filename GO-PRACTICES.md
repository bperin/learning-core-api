Core Principles for a "Not Overkill" Go DDD Style
Prioritize Strategic Design Principles: Focus on understanding the business domain, defining a shared "ubiquitous language" with domain experts, and identifying clear Bounded Contexts (conceptual boundaries). These high-level modeling techniques provide the most value with the least "complexity tax".
Use Go Idioms: Structure your code by feature or domain, not by rigid, deep layers inherited from other languages. Go favors simple, flat structures that remain maintainable.
Keep Dependencies Inverted: Your core domain logic should be independent of external concerns like databases or web frameworks. Use interfaces for decoupling (e.g., repository interfaces) and dependency injection via constructor functions.
Embrace Pragmatic Tactical Patterns: Apply tactical patterns like Entities, Value Objects, and Repositories when the specific part of your system has rich business rules.
Entities: Objects defined by their unique identity, not their attributes (e.g., a User with an ID).
Value Objects: Immutable objects defined by their attributes (e.g., a Money struct with amount and currency).
Repositories: Abstractions for data access, returning domain models, not database entities.
Avoid Dogma: Not every project needs advanced patterns like Event Sourcing or CQRS (Command Query Responsibility Segregation), which can add significant complexity. Start simple, and only add complexity when the business needs clearly justify it. Refactoring exists for a reason.
Recommended Code Structure
A common and effective approach in Go is using an architecture like Clean Architecture or Hexagonal Architecture (Ports and Adapters), which naturally align with DDD principles by keeping the domain isolated.
A typical project layout might look like this:
/internal/domain: Contains your core business logic, entities, value objects, domain events, and repository interfaces. This layer should have zero external dependencies.
/internal/application: Orchestrates domain services to perform use cases or application services.
/internal/infrastructure: Contains concrete implementations of interfaces defined in the domain layer, such as database repositories (using GORM, SQL, etc.), and interactions with third-party services.
/internal/interfaces (or /internal/delivery): Handles incoming requests (HTTP APIs, CLI commands) and maps them to application services.
