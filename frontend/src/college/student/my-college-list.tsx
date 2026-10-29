import type { ReactNode } from "react";
import { ArrowDown, ArrowUp } from "lucide-react";

import { errorMessage } from "../../api/client";
import type { PersonalCollegeApplication } from "../../api/models";
import { smallButtonClassName } from "../utils/shared";
import AddApplicationForm from "./add-application-form";
import ApplicationCard from "./application-card";
import DeadlineErrorDialog from "./deadline-error-dialog";
import { useMyCollegeList } from "./use-my-college-list";

// page frame with the "My List" heading, actions sit next to the heading
interface MyListShellProps {
  actions?: ReactNode;
  children: ReactNode;
}

function MyListShell({ actions, children }: MyListShellProps) {
  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-10">
      <div className="flex items-baseline gap-4">
        <h2 className="text-5xl font-medium">My List</h2>
        {actions}
      </div>
      {children}
    </div>
  );
}

export default function MyCollegeList() {
  const list = useMyCollegeList();

  if (list.isLoading) {
    return (
      <MyListShell>
        <p className="mt-6 text-sm">Loading...</p>
      </MyListShell>
    );
  }

  if (list.error) {
    return (
      <MyListShell>
        <p className="mt-6 text-sm">Error: {errorMessage(list.error)}</p>
      </MyListShell>
    );
  }

  const renderRankControls = (app: PersonalCollegeApplication) => {
    if (app.rank == null) {
      return (
        <button
          type="button"
          disabled={list.busy}
          onClick={() => list.addToRanking(app)}
          className={smallButtonClassName}
        >
          Add to ranking
        </button>
      );
    }

    const index = list.rankIndex(app);
    return (
      <div className="flex items-center gap-2">
        <button
          type="button"
          aria-label={`Move ${list.collegeName(app)} up`}
          disabled={list.busy || index === 0}
          onClick={() => list.moveUp(app)}
          className={smallButtonClassName}
        >
          <ArrowUp aria-hidden="true" className="h-3.5 w-3.5" />
        </button>
        <button
          type="button"
          aria-label={`Move ${list.collegeName(app)} down`}
          disabled={list.busy || index === list.rankedCount - 1}
          onClick={() => list.moveDown(app)}
          className={smallButtonClassName}
        >
          <ArrowDown aria-hidden="true" className="h-3.5 w-3.5" />
        </button>
        <button
          type="button"
          disabled={list.busy}
          onClick={() => list.unrank(app)}
          className={smallButtonClassName}
        >
          Unrank
        </button>
      </div>
    );
  };

  const renderCard = (app: PersonalCollegeApplication) => {
    const isEditing = list.editingId === app.id;

    return (
      <ApplicationCard
        key={app.id}
        app={app}
        college={list.collegeById.get(app.global_college_id)}
        collegeName={list.collegeName(app)}
        busy={list.busy}
        saving={list.saving}
        isEditing={isEditing}
        isDragging={list.draggedId === app.id}
        draggable={!list.busy && !isEditing}
        rankControls={renderRankControls(app)}
        onDragStart={() => list.setDraggedId(app.id)}
        onDragEnd={() => list.setDraggedId(null)}
        onDrop={() => list.handleDrop(app)}
        onStartEdit={() => list.startEditing(app)}
        onCancelEdit={list.cancelEditing}
        onSave={(applicationType, category) =>
          list.saveEdit(app, applicationType, category)
        }
        onRemove={() => list.removeApplication(app)}
      />
    );
  };

  const deadlineErrorCollege = list.deadlineError
    ? list.collegeById.get(list.deadlineError.globalCollegeId)
    : undefined;

  return (
    <MyListShell
      actions={
        <button
          type="button"
          onClick={() => list.setShowAddForm(!list.showAddForm)}
          className="text-sm underline"
        >
          {list.showAddForm ? "Close" : "Add a school"}
        </button>
      }
    >
      {list.showAddForm && (
        <AddApplicationForm
          colleges={list.availableColleges}
          initialCollegeId={list.preselectedId}
          busy={list.busy}
          adding={list.creating}
          onAdd={list.createApplication}
        />
      )}

      {list.mutationError && (
        <p role="alert" className="mt-3 text-sm text-red-600">
          Error: {errorMessage(list.mutationError)}
        </p>
      )}

      {list.applications.length === 0 ? (
        <p className="mt-6 text-sm">
          Your list is empty. Choose Add a school to start your list.
        </p>
      ) : (
        <>
          <p className="mt-4 text-xs text-gray-500">
            Drag colleges to reorder them, or use the arrow buttons. Drop a
            college onto an unranked one to unrank it.
          </p>
          <ol className="mt-4 flex flex-col gap-4">
            {list.applications.map(renderCard)}
          </ol>
        </>
      )}

      {list.deadlineError && (
        <DeadlineErrorDialog
          college={deadlineErrorCollege}
          collegeName={deadlineErrorCollege?.school_name ?? "This college"}
          applicationType={list.deadlineError.applicationType}
          onClose={list.closeDeadlineError}
        />
      )}
    </MyListShell>
  );
}
