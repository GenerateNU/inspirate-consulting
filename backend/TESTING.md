# Backend Testing Guide

This guide explains the testing strategy and conventions for the Go backend in Inspirate Consulting.

---

## Testing Strategy

We follow a 3-tier testing strategy:

| Level | Scope | Dependencies | When to Use | Location Example |
|---|---|---|---|---|
| **Handler Unit Tests** | Business logic & error mapping | None (Mocks) | Every handler feature / logic branch | `internal/handlers/greeting/create_test.go` |
| **Route / API Tests** | HTTP routing, status codes, Huma schema validation | None (Mocks) | Endpoint contract & validation testing | `internal/routes/greeting_test.go` |
| **Database Integration Tests** | SQL query syntax, constraints, and data mapping | Local PostgreSQL / Supabase | Verifying real SQL execution & queries | `internal/data/postgres/schema/greetingStore/create_test.go` |

---

## Mocking with Mockery

See [Mocks README](internal/data/repo-mocks/README.md) for more info.

---

## Running Tests

In `/backend` run all unit and fast tests:
```bash
go test -short ./...
```

Run all tests including integration tests:
```bash
go test ./...
```

Run a specific test:
```bash
go test -v ./internal/handlers/greeting -run TestHandler_CreateGreeting
```

Run tests with coverage:
```bash
go test -cover ./...
```

---

## Best Practices

1. **`t.Parallel()`**: Use `t.Parallel()` on top-level tests and subtests where state is isolated.
2. **Short Mode**: Mark any test requiring an external service (Docker / Supabase / PostgreSQL) with `if testing.Short() { t.Skip() }`.
3. **Mock Assertions**: Mockery automatically registers cleanup to assert expectations when initialized with `mocks.NewGreetingRepository(t)`.
