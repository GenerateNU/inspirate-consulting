# PostgreSQL data layer

This package is the concrete database implementation. It uses `pgxpool`, a
pool of reusable PostgreSQL connections, so concurrent HTTP requests do not
each open their own connection.

## Connection setup

[`db-connection.go`](db-connection.go) builds a connection string from the
database configuration, sets pool limits, creates the pool, and pings the
database before the server starts. If the connection cannot be established,
startup stops rather than serving an API that cannot reach its data.

The pool is created once and passed to the repository layer.

## Feature stores

SQL is grouped by feature under [`schema`](schema). For example,
[`greetingStore`](schema/greetingStore) implements the `GreetingRepository`
interface from [`internal/data/data.go`](../data.go). Its `CreateGreeting`
operation uses parameterized values (`$1`, `$2`) and the `greetings` table.

Use `QueryRow` for one returned row, `Query` for many rows, and `Exec` when no
rows are returned. Use `CollectRows` on pgx directly to scan and collect all rows 
from a table into a slice. Always pass the request context and parameterize user input;
never build SQL by concatenating strings.

See the [data layer guide](../README.md) for the interface and wiring steps,
and the [Supabase guide](../../supabase/README.md) for migrations.
