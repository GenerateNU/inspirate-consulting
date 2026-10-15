import { QueryClient } from "@tanstack/react-query";

import type { ErrorModel } from "./models";

// Shape of the errors thrown by the generated client on non-2xx responses.
export type ApiError = Error & { info?: ErrorModel; status?: number };

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // 4xx responses won't succeed on retry, so only retry network/5xx errors.
      retry: (failureCount, error: ApiError) =>
        !(error.status && error.status < 500) && failureCount < 3,
    },
  },
});

export const errorMessage = (error: unknown) => {
  const { info, message } = (error ?? {}) as ApiError;
  return info?.detail || info?.title || message || "Something went wrong.";
};
