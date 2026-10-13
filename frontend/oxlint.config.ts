import { defineConfig } from "oxlint";

export default defineConfig({
  ignorePatterns: [
    "src/api/endpoints/**",
    "src/api/models/**",
    "src/api/schemas/**",
  ],
  options: {
    typeAware: true,
    typeCheck: true,
  },
});
