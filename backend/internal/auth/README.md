# Authentication

This package talks to Supabase Auth for signup and login and verifies JWTs for
protected API routes.

`SupabaseSignup` and `SupabaseLogin` call Supabase's Auth API using the service
role key. Signup also checks for an eight-character password containing upper
and lower case letters, a digit, and a special character. These helpers are
not routes until they are registered in the routes package.

## Request authentication

The route setup installs `AuthMiddleware` before protected routes unless
`TEST_MODE=true`. The middleware expects a cookie named `jwt`, verifies it with
the `SUPABASE_JWT_SECRET` environment variable, and returns `401` when the
cookie is missing or invalid. It does not currently read a Bearer token from an
`Authorization` header.

Keep public route registration before the middleware and protected route
registration after it. The current signup and login route registration is
commented out in [`routes/server.go`](../routes/server.go), so the helpers are
available but not exposed by the API yet.