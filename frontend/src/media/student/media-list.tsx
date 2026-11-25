import { useState } from "react";
import { keepPreviousData } from "@tanstack/react-query";

import { errorMessage } from "../../api/client";
import { useGetMediaAccess } from "../../api/endpoints";
import { useDebouncedSearch } from "../shared/hooks/use-debounced-search";
import ModuleCard from "../shared/module-card";
import ModuleListFrame from "../shared/module-list-frame";
import ModuleListSkeleton from "../shared/module-list-skeleton";
import ModuleSearch from "../shared/module-search";
import ProgressCard from "./progress-card";
import StatusBadge from "./status-badge";
import { formatDueDate } from "./utils/format-due-date";
import { getModuleStatus } from "./utils/module-status";

// student's "Modules" page: the modules they've been assigned, with their progress
export default function StudentMediaList() {
  const [query, setQuery] = useState("");
  const debouncedQuery = useDebouncedSearch(query);

  // 100 is the backend's max page size
  const allModules = useGetMediaAccess({ limit: 100 });
  // same request as allModules (so the same cache entry) until a search is typed
  const shownModules = useGetMediaAccess(
    { q: debouncedQuery || undefined, limit: 100 },
    { query: { placeholderData: keepPreviousData } },
  );

  const all = allModules.data ?? [];
  const shown = shownModules.data ?? [];
  const completedCount = all.filter((item) => item.completed_at).length;

  let body;
  if (shownModules.isLoading) {
    body = <ModuleListSkeleton />;
  } else if (shownModules.error) {
    body = (
      <p className="font-host text-sm">
        Error: {errorMessage(shownModules.error)}
      </p>
    );
  } else if (shown.length === 0) {
    body = (
      <p className="font-host text-sm text-dark-greyish-brown">
        {debouncedQuery
          ? `No modules match "${debouncedQuery}".`
          : "No modules have been assigned to you yet."}
      </p>
    );
  } else {
    body = (
      <div role="list" className="flex w-full flex-col gap-4">
        {shown.map((item) => {
          const status = getModuleStatus(item);
          return (
            <ModuleCard
              key={item.id}
              to={`/student-media-list/${item.id}`}
              title={item.title}
              description={item.description}
              muted={status === "completed"}
              details={
                <>
                  <StatusBadge status={status} />
                  {item.due_date && (
                    <span className="shrink-0 font-host text-xs text-dark-greyish-brown">
                      Due: {formatDueDate(item.due_date)}
                    </span>
                  )}
                </>
              }
            />
          );
        })}
      </div>
    );
  }

  return (
    <ModuleListFrame
      actions={<ModuleSearch value={query} onChange={setQuery} />}
    >
      {allModules.data && (
        <ProgressCard completed={completedCount} total={all.length} />
      )}
      {body}
    </ModuleListFrame>
  );
}
