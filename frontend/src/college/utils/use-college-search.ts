import { useState } from "react";
import { keepPreviousData } from "@tanstack/react-query";

import {
  useListGlobalColleges,
  useSearchGlobalColleges,
} from "../../api/endpoints";
import type { GlobalCollege } from "../../api/models";
import { useDebouncedSearch } from "../../media/shared/hooks/use-debounced-search";

// shared state for the college list pages:
// shows every college until the user types a query, then shows the search results instead
export function useCollegeSearch() {
  const [query, setQuery] = useState("");
  // only search once typing pauses, not on every keystroke
  const debouncedQuery = useDebouncedSearch(query);
  const searching = debouncedQuery !== "";

  const list = useListGlobalColleges();

  const search = useSearchGlobalColleges(
    { query: debouncedQuery },
    // keep showing the last results while the next search loads, so the list doesn't flash "Loading..."
    { query: { enabled: searching, placeholderData: keepPreviousData } },
  );

  const colleges: GlobalCollege[] = (searching ? search.data : list.data) ?? [];

  // refetch both the full list and the current search so neither goes stale
  const refetch = async () => {
    await list.refetch();
    if (searching) await search.refetch();
  };

  return {
    query,
    setQuery,
    searching,
    colleges,
    isLoading: list.isLoading,
    error: list.error,
    searchLoading: search.isLoading,
    searchError: search.error,
    refetch,
  };
}

export type CollegeSearchState = ReturnType<typeof useCollegeSearch>;
