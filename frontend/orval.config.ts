import { defineConfig } from "orval";

export default defineConfig({
  // HTTP client generation
  clients: {
    input: {
      target: "../backend/internal/api/openapi.yaml",
    },
    output: {
      mode: "tags-split",
      client: "react-query",
      httpClient: "fetch",
      target: "src/api/endpoints",
      schemas: "src/api/models",
      mock: true,
      // Relative so the same bundle works on localhost, previews and prod.
      // The '/api' prefix is stripped by the Vite dev proxy in development and
      // by the Netlify proxy in production, before the backend sees it.
      baseUrl: "/api",
      override: {
        fetch: {
          // Return the response body directly and throw on non-2xx, so
          // TanStack Query surfaces failures through its `error` state.
          includeHttpResponseReturnType: false,
          forceSuccessResponse: true,
        },
      },
    },
  },
  zodSchemas: {
    output: {
      client: "zod",
      mode: "single",
      target: "./src/api/schemas",
      override: {
        zod: {
          // Prefer Mini for more tree-shakeable generated schemas.
          variant: "mini",
          version: 4,
        },
      },
    },
    input: {
      // Path or URL to our openapi spec
      target: "../backend/internal/api/openapi.yaml",
    },
  },
});
