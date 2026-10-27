import type { ReactNode } from "react";

import { errorMessage } from "../api/client";
import type { GlobalCollege } from "../api/models";
import { SearchIcon } from "./shared";
import type { CollegeSearchState } from "./use-college-search";

// search box plus the list of matching colleges
// each page decides how a card renders and what goes in the footer
export default function CollegeSearch({
  title,
  search,
  renderCollege,
  notice,
  footer,
}: {
  title: string;
  search: CollegeSearchState;
  renderCollege: (college: GlobalCollege) => ReactNode;
  notice?: ReactNode;
  footer?: ReactNode;
}) {
  const { query, setQuery, searching, colleges, searchLoading, searchError } =
    search;

  let body;
  if (search.isLoading) {
    body = <p className="py-6 text-center text-sm">Loading...</p>;
  } else if (search.error) {
    body = <p className="py-6 text-center text-sm">Something went wrong.</p>;
  } else {
    body = (
      <>
        {searching && searchLoading && (
          <p className="py-6 text-center text-sm">Searching...</p>
        )}
        {searching && searchError && (
          <p className="py-6 text-center text-sm">
            Error: {errorMessage(searchError)}
          </p>
        )}
        {!searchLoading && colleges.length === 0 && (
          <p className="py-6 text-center text-sm">
            {searching ? "No matching colleges." : "No colleges yet."}
          </p>
        )}
        <ul className="flex flex-col gap-4">
          {colleges.map((college) => renderCollege(college))}
        </ul>
      </>
    );
  }

  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-10">
      <h2 className="mb-4 text-2xl text-gray-500">{title}</h2>
      <div className="flex min-h-[32rem] flex-col border border-black p-4">
        <label className="flex items-center gap-3 border-b border-gray-300 px-2 pb-3">
          <SearchIcon />
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search by name or location"
            aria-label="Search colleges"
            className="w-full bg-transparent text-lg font-semibold outline-none placeholder:font-normal placeholder:text-gray-400"
          />
        </label>

        {notice}

        <div className="mt-4 flex-1">{body}</div>

        {footer && (
          <div className="mt-6 flex flex-col items-center gap-3 border-t border-gray-500 pt-6">
            {footer}
          </div>
        )}
      </div>
    </div>
  );
}
