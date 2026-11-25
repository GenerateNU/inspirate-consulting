import { useState } from "react";
import { Link } from "react-router-dom";
import { keepPreviousData } from "@tanstack/react-query";
import { Plus, Trash2 } from "lucide-react";

import { errorMessage } from "../../api/client";
import { useDeleteMedia, useListMedia } from "../../api/endpoints";
import type { Media } from "../../api/models";
import { Button } from "@/components/ui/button";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { cn } from "@/lib/utils";
import { useDebouncedSearch } from "../shared/hooks/use-debounced-search";
import { primaryButtonClassName } from "../shared/media-styles";
import ModuleCard from "../shared/module-card";
import ModuleListFrame from "../shared/module-list-frame";
import ModuleListSkeleton from "../shared/module-list-skeleton";
import ModuleSearch from "../shared/module-search";
import { ErrorMessage } from "@/components/ui/error-message";

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
    body = <ModuleListSkeleton />;
  } else if (mediaQuery.error) {
    body = (
      <ErrorMessage
        variant="banner"
        title="Couldn't load modules"
        error={mediaQuery.error}
        onRetry={mediaQuery.refetch}
      />
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
          <ModuleCard
            key={item.id}
            to={`/counselor-media-list/${item.id}`}
            title={item.title}
            description={item.description}
            action={
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
            }
          />
        ))}
      </div>
    );
  }

  return (
    <ModuleListFrame
      actions={
        <>
          <ModuleSearch value={query} onChange={setQuery} />
          <Button asChild className={cn(primaryButtonClassName, "px-6")}>
            <Link to="/counselor-media-list/new">
              <Plus strokeWidth={2} />
              Create Module
            </Link>
          </Button>
        </>
      }
    >
      {body}

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
    </ModuleListFrame>
  );
}
