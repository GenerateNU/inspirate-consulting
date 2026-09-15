# Application entry point

[`main.go`](main.go) is the executable package. Run it from the `backend`
directory with:

```sh
go run ./cmd
```

It loads environment variables into `config.Config`, initializes the routes and
database repository, starts Fiber on `PORT` (8080 by default), and waits for a
shutdown signal. `Ctrl+C` or a deployment stop signal lets the server close
HTTP connections and the database pool cleanly.

Go does not read `.env` automatically. Export the file in your shell first, as
shown in the [setup guide](../../SETUP.md), or provide variables through your
container or hosting platform.

`TEST_MODE=true` skips authentication middleware for local testing. Do not use
it in a deployed environment.