import { useState } from "react";

import {
  useListGlobalColleges,
  useSearchGlobalColleges,
} from "../../api/endpoints";
import type { GlobalCollege } from "../../api/models";

// shared state for the college list pages:
// shows every college until the user types a query, then shows the search results instead
export function useCollegeSearch() {
  const [query, setQuery] = useState("");
  const trimmedQuery = query.trim();
  const searching = trimmedQuery !== "";

  const list = useListGlobalColleges();

  const search = useSearchGlobalColleges(
    { query: trimmedQuery },
    { query: { enabled: searching } },
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
