# Routes

Routes connect an HTTP path and method to a handler. This package also creates
the Fiber server, installs middleware, and adapts Fiber to Huma.

## Startup flow

[`server.go`](server.go) builds the application in this order:

1. Create a Fiber app with the central HTTP error handler.
2. Install recovery, compression, logging, CORS, and favicon middleware.
3. Create a Huma API on top of Fiber. Huma validates inputs and generates
	 OpenAPI documentation.
4. Install authentication middleware unless `TEST_MODE=true`.
5. Register the root route and protected feature routes.

Middleware only affects routes registered after it. Keep public routes before
the auth middleware and protected routes after it.

The greeting route is registered in [`greeting.go`](greeting.go). Add new route
registration there only when it belongs to the greeting feature; otherwise
create a feature-specific file and call it from `setupProtectedHumaRoutes`.

See the [handler guide](../handlers/README.md) for what should happen after a
route receives a request and the [models guide](../models/README.md) for Huma
input and output structs.