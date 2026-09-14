import { defineConfig } from 'orval';

export default defineConfig({
  // HTTP client generation
  clients: {
    input: {
      target: '../backend/internal/api/openapi.yaml',
    },
    output: {
      mode: 'tags-split',
      client: 'swr',
      target: 'src/api/endpoints',
      schemas: 'src/api/models',
      mock: true,
      baseUrl: 'http://127.0.0.1:3001'
    },
  },
  zodSchemas: {
    output: {
      client: 'zod',
      mode: 'single',
      target: './src/api/schemas',
      override: {
        zod: {
          // Prefer Mini for more tree-shakeable generated schemas.
          variant: 'mini',
          version: 4,
        },
      },
    },
    input: {
      // Path or URL to our openapi spec
      target: '../backend/internal/api/openapi.yaml',
    },
  },
});