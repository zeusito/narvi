# AI Agents

This document provides context and guidelines for AI agents working on the Salus project.

## Project Overview

This is our main API for our platform. We aim for code that is highly maintainable, testable, and scalable. Stability is our priority.

- **Frameworks**: Chi (Router), Bun (ORM), Zerolog (Logging).
- **Structure**:
  - `cmd/`: Application entry points.
  - `internal/`: Core business logic.
  - `pkg/`: Shared infrastructure and tools.

**File Organization:**
We follow a strict separation of concerns. Each module should have its own directory under `internal/`. A module is a self-contained unit of code that handles a specific feature.

| File                 | Purpose                                                                                                                                                                                  |
| :------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `factory.go`         | **Entry Point.** Contains `NewModule(...)`. Wires the Repository, Service, and Controller together.                                                                                      |
| `controller.go`      | **Transport Layer.** Contains the `Controller` struct, HTTP handlers, request binding, validation, and response rendering. Registers routes with the `chi` router. Depends on `Service`. |
| `service.go`         | **Business Logic.** Defines the `Service` interface and `NewService` constructor.                                                                                                        |
| `default_service.go` | **Logic Implementation.** Contains `DefaultService` (or `ServiceDefault`) which implements `Service`.                                                                                    |
| `repository.go`      | **Data Access.** Defines the `Repository` interface and `NewRepository` constructor.                                                                                                     |
| `default_repo.go`    | **DB Implementation.** Contains `DefaultRepository` (or `RepositoryDefault`) which implements `Repository` (typically using `bun`).                                                      |
| `models.go`          | **Data Structures.** Contains domain models, DTOs (Request/Response structs), and value objects specific to this module.                                                                 |

## Coding Conventions

- **Dependency Injection:** Pass dependencies (DB, other services) into the `NewModule` or `NewService` constructors.
- **Router:** We use `chi`. Controllers should accept `*chi.Mux` (or `chi.Router`) and register their own sub-routes.
- **Database:** We use `uptrace/bun`. Repositories should accept `*bun.DB` (or `bun.IDB`).
- **Error Handling:** Use `pkg/terrors` for typed errors in the service/controller layer.
- **Idiomatic Go**: Follow standard Go practices and patterns.
- **Separation of Concerns**: Keep business logic in `internal/` and infrastructure in `pkg/`.
- **Consistency**: Match the existing coding style in the repository.
- **Automation**: Use the `Makefile` for tasks like building, running, and testing.

## Testing Rules

- **Required for Go changes:** run `make test` (or `go test -v ./... --race`) before handing off work.
- **Scoped changes:** also run tests for the modified package(s) when feasible (e.g., `go test ./internal/wallets`).
- **DB-dependent tests:** if you add or modify tests that require Postgres, document the setup in the PR or task notes and ensure migrations are up to date.
- **Test style:** avoid table-driven tests; prefer assertion-style tests.
- **Avoid mocking:** prefer testing against a real Postgres instance (e.g., using `testcontainers`) rather than mocking the database layer. If mocking is necessary, use interfaces and dependency injection to facilitate testing.

## Task Execution

Before starting a task:

1. Review `README.md` for project context.
2. Check existing patterns in `internal/iam` or `internal/healthcheck`.
3. Verify changes using `make test` and `make lint`.
