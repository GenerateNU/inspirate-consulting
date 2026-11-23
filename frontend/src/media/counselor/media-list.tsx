import { useState } from "react";
import { Link } from "react-router-dom";
import { keepPreviousData } from "@tanstack/react-query";
import { Plus, Search, SquarePlay, Trash2 } from "lucide-react";

import { errorMessage } from "../../api/client";
import { useDeleteMedia, useListMedia } from "../../api/endpoints";
import type { Media } from "../../api/models";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";
import { useDebouncedSearch } from "./hooks/use-debounced-search";
import { panelClassName, primaryButtonClassName } from "./media-styles";

// counselor's "Modules" page: the whole media library, searchable
export default function CounselorMediaList() {
  const [query, setQuery] = useState("");
  const debouncedQuery = useDebouncedSearch(query);
  const [pendingDelete, setPendingDelete] = useState<Media | null>(null);

  // 100 is the backend's max page size
  const mediaQuery = useListMedia(
    { q: debouncedQuery || undefined, limit: 100 },
    { query: { placeholderData: keepPreviousData } },
  );
  const deleteMedia = useDeleteMedia();
  const media = mediaQuery.data ?? [];

  const confirmDelete = () => {
    if (!pendingDelete) return;
    deleteMedia.mutate(
      { id: pendingDelete.id },
      {
        onSuccess: async () => {
          await mediaQuery.refetch();
          setPendingDelete(null);
        },
      },
    );
  };

  let body;
  if (mediaQuery.isLoading) {
    body = (
      <div className="flex w-full flex-col gap-4">
        {[0, 1, 2].map((i) => (
          <Skeleton key={i} className="h-[115px] w-full rounded-lg" />
        ))}
      </div>
    );
  } else if (mediaQuery.error) {
    body = (
      <p className="font-host text-sm">
        Error: {errorMessage(mediaQuery.error)}
      </p>
    );
  } else if (media.length === 0) {
    body = (
      <p className="font-host text-sm text-dark-greyish-brown">
        {debouncedQuery
          ? `No modules match "${debouncedQuery}".`
          : "No modules yet. Choose Create Module to add one."}
      </p>
    );
  } else {
    body = (
      <div role="list" className="flex w-full flex-col gap-4">
        {media.map((item) => (
          <Card
            key={item.id}
            role="listitem"
            className={cn(
              panelClassName,
              "relative min-h-[115px] flex-row items-center gap-4 overflow-clip py-6 pr-6 pl-7",
            )}
          >
            <div
              aria-hidden="true"
              className="absolute top-[-1px] bottom-[-1px] left-[-1px] w-1 bg-tan"
            />
            <div className="flex size-9 shrink-0 items-center justify-center rounded-lg bg-lightest-yellow">
              <SquarePlay
                className="size-5 text-dark-greyish-brown"
                strokeWidth={1.75}
              />
            </div>
            {/* inline-size containment: long titles/descriptions no longer
                widen the page; they truncate to fit the card instead */}
            <div className="flex min-w-0 flex-1 flex-col gap-1 [contain:inline-size]">
              <Link
                to={`/counselor-media-list/${item.id}`}
                className="truncate font-jakarta text-lg leading-[normal] font-medium text-charcoal hover:underline"
              >
                {item.title}
              </Link>
              <p className="truncate font-host text-xs text-dark-greyish-brown">
                {item.description}
              </p>
            </div>
            <Button
              type="button"
              variant="ghost"
              size="icon"
              aria-label={`Delete ${item.title}`}
              title="Delete module"
              onClick={() => {
                deleteMedia.reset();
                setPendingDelete(item);
              }}
              className="rounded-lg text-dark-greyish-brown hover:bg-off-white"
            >
              <Trash2 className="size-[18px]" strokeWidth={1.75} />
            </Button>
          </Card>
        ))}
      </div>
    );
  }

  return (
    <div className="flex min-h-full w-full flex-col gap-6 bg-off-white py-12 pr-20 pl-[58px] max-lg:px-6">
      <div className="flex w-full flex-wrap items-center gap-4 pl-[59px] max-lg:pl-0">
        <h1 className="min-w-0 flex-1 font-jakarta text-2xl leading-[normal] font-bold text-charcoal">
          Modules
        </h1>
        <div className="relative w-[264px] max-md:w-full">
          <Search
            aria-hidden="true"
            className="pointer-events-none absolute top-1/2 left-3 size-3.5 -translate-y-1/2 text-warm-grey"
            strokeWidth={2}
          />
          <Input
            type="search"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search specific modules"
            aria-label="Search modules"
            className="h-8 rounded-full border-cream bg-white pr-3 pl-8 font-host text-xs text-charcoal shadow-none placeholder:text-warm-grey md:text-xs"
          />
        </div>
        <Button asChild className={cn(primaryButtonClassName, "px-6")}>
          <Link to="/counselor-media-list/new">
            <Plus strokeWidth={2} />
            Create Module
          </Link>
        </Button>
      </div>
      <div className="flex w-full flex-col pl-[57px] max-lg:pl-0">{body}</div>

      <ConfirmDialog
        open={pendingDelete != null}
        onOpenChange={(open) => !open && setPendingDelete(null)}
        title="Are you sure you want to delete this module?"
        description={
          deleteMedia.error ? (
            <span className="text-danger">
              Error: {errorMessage(deleteMedia.error)}
            </span>
          ) : (
            "Once deleted, this cannot be recovered."
          )
        }
        confirmLabel="Delete Module"
        pendingLabel="Deleting..."
        pending={deleteMedia.isPending}
        onConfirm={confirmDelete}
      />
    </div>
  );
}
