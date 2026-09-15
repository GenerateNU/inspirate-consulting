# Configuration

Configuration is read from environment variables when the application starts.
`envconfig` maps names such as `DB_HOST` into the structs in this package. The
Go process does not load `.env` files itself.

## Main variables

- `PORT` defaults to `8080`.
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, and `DB_NAME` are required.
- `DB_SSLMODE` defaults to `require`; use `disable` for local PostgreSQL when
  appropriate.
- `SUPABASE_URL`, `SUPABASE_ANON_KEY`, and `SUPABASE_SERVICE_ROLE_KEY` are
  required by the configuration.
- `TEST_MODE=true` is read separately and disables route authentication for
  local testing only.

Database pool limits and application URLs are also defined in
[`db.go`](db.go) and [`application.go`](application.go). Add a field here when
new configuration is needed instead of reading environment variables throughout
the application.