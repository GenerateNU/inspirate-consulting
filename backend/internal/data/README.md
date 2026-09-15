# Data Layer

The data layer is the only part of the application that should know how data is
stored. Handlers ask it to perform operations such as "create a greeting"; the
PostgreSQL implementation decides which SQL to run.

Keeping SQL here prevents routes and handlers from depending on database
details. It also lets handler tests replace PostgreSQL with a mock.

## Request Flow

For the greeting feature, a database call moves through these pieces:

```text
Greeting handler
	-> GreetingRepository interface in data.go
	-> PostgreSQL implementation in postgres/schema/greetingStore
	-> pgx connection pool
	-> PostgreSQL
```

The caller passes a `context.Context` through every step. The context allows a
database query to be cancelled if the HTTP request ends or times out.

## Interfaces and Implementations

[`data.go`](data.go) defines repository interfaces. An interface lists the
methods a handler is allowed to call.

An interface contains no SQL. A concrete type satisfies it by implementing all
of its methods. Go does this automatically; there is no `implements` keyword.

The top-level `Repository` groups the feature repositories used by the app.
`NewRepository` receives one database pool and constructs each concrete
repository.

## Folder Guide

- [`data.go`](data.go) defines repository interfaces and wires implementations.
- [`postgres/db-connection.go`](postgres/db-connection.go) opens and verifies
	the shared PostgreSQL connection pool.
- [`postgres/schema`](postgres/schema) contains SQL grouped by feature.
- [`repo-mocks`](repo-mocks/README.md) contains test replacements for repository
	interfaces.
- [`db-interface`](db-interface/db-interface.go) defines query methods shared by
	pools and transactions.

See the [PostgreSQL guide](postgres/README.md) for query-writing conventions and
the [Supabase guide](../supabase/README.md) for tables and migrations.

## pgx and pgxpool

**pgx** is the Go library used to send parameterized SQL to PostgreSQL and scan
returned columns into Go values.

**pgxpool** manages several reusable database connections. A web server may
handle many requests at once, while a single connection can process only one
query at a time. The pool lends a connection to each query and returns it when
the query finishes. Application code should share the pool instead of opening a
new connection for every request.

Common pgx methods are:

- `QueryRow(...).Scan(...)` for an operation returning one row.
- `Query(...)` for an operation returning multiple rows.
- `Exec(...)` for an operation returning no rows.

Always use placeholders such as `$1` and `$2` for values. Do not join user input
into an SQL string.

## Adding a Data Operation

For another operation on an existing feature:

1. Add the method signature to the feature interface in `data.go`.
2. Add an operation file, such as `get.go`, to its PostgreSQL store package.
3. Implement the method with the exact same parameters and return types.
4. Update or regenerate the repository mock.
5. Add tests and run `go test ./...` from `backend`.

For a new feature, also create its store package, add it as a field on
`Repository`, and initialize it in `NewRepository`. If its table or columns are
new, create a [Supabase migration](../supabase/README.md) before running the query.

Repository methods must return database errors to their caller,
the handler and central HTTP error layer need the error to produce the correct response 
and log useful details.