# Handlers

A handler contains the application action behind an endpoint. It receives
validated input from Huma, calls a repository, and returns a model or an
error. It should not know how SQL is written.

## Greeting example

[`greeting`](greeting) defines a `Handler` with a
`GreetingRepository` dependency. `NewHandler` receives that dependency, so the
real PostgreSQL repository can be replaced by a mock in a test. `CreateGreeting`
passes the request context and input to the repository and returns its result.

The context must travel through every layer. It lets a database call stop when
the request is cancelled or times out.

Handlers are intentionally thin at the moment. Put feature decisions here;
put persistence details in the repository and request/response shapes in
[`models`](../models/README.md).

## Adding a handler

Create a feature package, define a handler dependency on a data interface, and
construct it in the feature's route setup. Return repository errors unchanged
so the central error handler can choose the HTTP response. Use
[`repo-mocks`](../data/repo-mocks/README.md) to test the handler without a
database.
