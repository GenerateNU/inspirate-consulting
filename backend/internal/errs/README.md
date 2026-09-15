# HTTP errors

[`http.go`](http.go) defines the error format returned by the API:

```json
{"code": 400, "message": "invalid request"}
```

Use the helpers such as `BadRequest`, `Unauthorized`, `NotFound`, and
`InternalServerError` when code needs to return a known HTTP failure. Returning
an `HTTPError` preserves its status code and message.

Fiber uses `ErrorHandler` as the central fallback. Unknown errors are logged
and returned as a generic `500` response so internal details are not sent to
clients. The handler also records the HTTP method and path in the server log.

Handlers should return errors rather than writing responses themselves. This
keeps response formatting consistent across routes.