# Models

Models are Go structs that describe data crossing an API boundary. Huma uses
them to validate requests and generate the OpenAPI document; the frontend can
use that contract to generate TypeScript types.

Tags on a field are instructions for Huma, for example in the greeting models:

- `path:"name"` reads a URL path value instead of a JSON body.
- `maxLength:"30"` rejects names longer than the database column allows.
- `example` and `doc` improve generated API documentation.
- `json:"greeting"` controls the response JSON name.

When changing a model, check the generated API contract and the frontend
callers. Keep models focused on transport data; database-specific code belongs
in the data layer.
