import { useEffect, useState } from "react";

// trimmed copy of the search text that only updates once typing pauses so the list isn't refetched on every keystroke
export function useDebouncedSearch(query: string, delayMs = 300) {
  const [debouncedQuery, setDebouncedQuery] = useState("");

  useEffect(() => {
    const timeout = setTimeout(() => setDebouncedQuery(query.trim()), delayMs);
    return () => clearTimeout(timeout);
  }, [query, delayMs]);

  return debouncedQuery;
}
