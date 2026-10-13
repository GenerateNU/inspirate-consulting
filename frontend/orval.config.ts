import { defineConfig } from "orval";

export default defineConfig({
  // HTTP client generation
  clients: {
    input: {
      target: "../backend/internal/api/openapi.yaml",
    },
    output: {
      mode: "tags-split",
      client: "swr",
      target: "src/api/endpoints",
      schemas: "src/api/models",
      mock: true,
      // Relative so the same bundle works on localhost, previews and prod.
      // The '/api' prefix is stripped by the Vite dev proxy in development and
      // by the Netlify proxy in production, before the backend sees it.
      baseUrl: "/api",
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
