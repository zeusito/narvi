---
name: module-pattern
description: Guidelines and templates for creating and testing standard internal modules in this project.
---

# Module Pattern Skill

Use this skill when creating a new module or adding functionality to an existing one in `internal/`. This ensures strict adherence to the project architecture: Chi (Router) -> Service -> Repository -> Bun (ORM).

## Module Architecture

Every module must be structured as follows:

| File                 | Purpose                              | Patterns to Follow                                        |
| :------------------- | :----------------------------------- | :-------------------------------------------------------- |
| `factory.go`         | **Entry Point.** Wires dependencies. | Use `NewModule(db *bun.DB, router chi.Router)`            |
| `controller.go`      | **Transport Layer.** HTTP handlers.  | Use `chi` for routing, `render` for responses.            |
| `service.go`         | **Business Logic Interface.**        | Define the interface here.                                |
| `service_default.go` | **Logic Implementation.**            | Implement the `Service` interface. Use `PrincipalClaims`. |
| `repository.go`      | **Data Access Interface.**           | Define the interface here.                                |
| `repo_default.go`    | **DB Implementation.**               | Use `uptrace/bun`. Implementation should be isolated.     |
| `models.go`          | **Data Structures.**                 | Domain models (Bun tags), Request/Response DTOs, Filters. |

## Implementation Rules

1.  **Dependency Injection**: Dependencies like `*bun.DB` or other services must be passed via constructors.
2.  **Context & Logging**: Every Service/Repo method must accept `context.Context`. Use `toolbox.GetRequestID(ctx)` and `claims.PrincipalID` in every log line.
3.  **Error Handling**: Use `pkg/terrors` (e.g., `terrors.RecordNotFound`, `terrors.OperationFailed`). NEVER return raw DB errors from the service layer.
4.  **Pagination & Filters**: Controllers should use a `Filters` struct from `models.go` and pass it to the service.

## Testing Guidelines

Ares prioritizes integration tests over mocking. Tests must use a real Postgres container.

### 1. TestMain Setup

Use `pkg/testbox` to initialize a container. Include `../../db/schema.sql` and a module-specific seed file.

```go
func TestMain(m *testing.M) {
    initScriptPath := []string{
        "../../db/schema.sql",
        "../../db/testdata/your_module.sql",
    }
    connData, closeFunc, err := testbox.InitPostgresqlContainer(context.Background(), initScriptPath)
    // ... setup testDB ...
    m.Run()
}
```

### 2. Test Style

- Use `github.com/stretchr/testify/assert`.
- Use `t.Context()` for the context.
- Avoid table-driven tests; prefer detailed assertion blocks for readability.
- Always include a "Seed" file in `db/testdata/` for consistency.
