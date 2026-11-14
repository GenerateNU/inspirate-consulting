import { useState, type ReactNode } from "react";

import { errorMessage } from "../../api/client";
import arrowsDownUpIcon from "@/assets/college/arrows-down-up.svg";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuLabel,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import AddSchoolModal from "./add-school-modal";
import ApplicationCard from "./application-card";
import RemoveSchoolDialog from "./remove-school-dialog";
import {
  sortOptions,
  useMyCollegeList,
  type SortKey,
} from "./use-my-college-list";

// page frame with the "My List" heading; actions sit to the right of it
function MyListShell({
  actions,
  children,
}: {
  actions?: ReactNode;
  children: ReactNode;
}) {
  return (
    <div className="flex min-h-full w-full flex-col gap-5 bg-off-white py-12 pr-20 pl-[115px] max-lg:px-6">
      <div className="flex w-full items-center justify-between gap-4">
        <h1 className="font-jakarta text-[34px] leading-[normal] font-bold tracking-[-0.34px] text-black">
          My List
        </h1>
        {actions}
      </div>
      {children}
    </div>
  );
}

export default function MyCollegeList() {
  const list = useMyCollegeList();
  const [addOpen, setAddOpen] = useState(false);

  if (!list.isStudent) {
    return (
      <MyListShell>
        <p className="font-host text-sm text-dark-greyish-brown">
          My List is for students. Log in as a student to see your list.
        </p>
      </MyListShell>
    );
  }

  if (list.isLoading) {
    return (
      <MyListShell>
        <p className="font-host text-sm">Loading...</p>
      </MyListShell>
    );
  }

  if (list.error) {
    return (
      <MyListShell>
        <p className="font-host text-sm">Error: {errorMessage(list.error)}</p>
      </MyListShell>
    );
  }

  const sortLabel = sortOptions.find(
    ({ value }) => value === list.sortKey,
  )?.label;

  return (
    <MyListShell
      actions={
        <button
          type="button"
          data-slot="primary-action"
          onClick={() => setAddOpen(true)}
          className="rounded-lg bg-dark-green px-4 py-2 font-host text-xs leading-[normal] text-white hover:bg-dark-green/90"
        >
          + Add a School
        </button>
      }
    >
      <div className="-mt-1 flex w-full items-center justify-end">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button
              type="button"
              data-slot="sort-trigger"
              aria-label={`Sort by: ${sortLabel}`}
              className="flex items-center justify-center gap-px rounded-[15px] border border-greyish-brown bg-white px-2.5 py-2 font-host text-xs leading-[normal] text-charcoal hover:bg-off-white"
            >
              <img src={arrowsDownUpIcon} alt="" width={16} height={16} />
              Sort By
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="font-host">
            <DropdownMenuLabel>Sort by</DropdownMenuLabel>
            <DropdownMenuRadioGroup
              value={list.sortKey}
              onValueChange={(value) => list.setSortKey(value as SortKey)}
            >
              {sortOptions.map(({ value, label }) => (
                <DropdownMenuRadioItem key={value} value={value}>
                  {label}
                </DropdownMenuRadioItem>
              ))}
            </DropdownMenuRadioGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {list.applications.length === 0 ? (
        <p className="font-host text-sm text-dark-greyish-brown">
          Your list is empty. Choose + Add a School to start your list.
        </p>
      ) : (
        <ol className="flex w-full flex-col gap-5">
          {list.applications.map((app) => (
            <ApplicationCard
              key={app.id}
              app={app}
              college={list.collegeById.get(app.global_college_id)}
              collegeName={list.collegeName(app)}
              onRemove={() => list.requestRemoval(app)}
            />
          ))}
        </ol>
      )}

      <AddSchoolModal
        open={addOpen}
        onOpenChange={setAddOpen}
        addedCollegeIds={list.addedCollegeIds}
        adding={list.adding}
        error={list.addError}
        onSave={list.addSchools}
      />

      <RemoveSchoolDialog
        collegeName={
          list.pendingRemoval ? list.collegeName(list.pendingRemoval) : null
        }
        removing={list.removing}
        error={list.removeError}
        onCancel={list.cancelRemoval}
        onConfirm={list.confirmRemoval}
      />
    </MyListShell>
  );
}
