import { useState, type ReactNode } from "react";
import { Link, useSearchParams } from "react-router-dom";

import { errorMessage, type ApiError } from "../../api/client";
import {
  useListPersonalCollegeApplications,
  useListGlobalColleges,
  useCreatePersonalCollegeApplication,
  useUpdatePersonalCollegeApplication,
  useUpdatePersonalCollegeApplicationRank,
  useDeletePersonalCollegeApplication,
} from "../../api/endpoints";
import {
  CreatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType,
  CreatePersonalCollegeApplicationRequestBodyCategory as Category,
  type GlobalCollege,
  type PersonalCollegeApplication,
} from "../../api/models";
import { deadlineOptions } from "../shared";
import AddApplicationForm from "./add-application-form";
import ApplicationCard, { smallButtonClassName } from "./application-card";
import DeadlineErrorDialog from "./deadline-error-dialog";

// page frame with the "My List" heading, actions sit next to the heading
function MyListShell({
  actions,
  children,
}: {
  actions?: ReactNode;
  children: ReactNode;
}) {
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
  const [draggedId, setDraggedId] = useState<number | null>(null);
  const [editingId, setEditingId] = useState<number | null>(null);

  // set when the backend rejects a deadline the college doesn't offer -> shows the popup
  const [deadlineError, setDeadlineError] = useState<{
    globalCollegeId: number;
    applicationType: string;
  } | null>(null);

  // ?college=<id> preselects a college from the colleges list and opens the form
  const [searchParams] = useSearchParams();
  const [showAddForm, setShowAddForm] = useState(searchParams.has("college"));

  const {
    data: applicationsData,
    error,
    isLoading,
    refetch,
  } = useListPersonalCollegeApplications();

  const { data: collegesData, isLoading: collegesLoading } = useListGlobalColleges();

  // one hook per kind of change, the row's id is passed when calling mutate
  const createApplication = useCreatePersonalCollegeApplication();
  const updateApplication = useUpdatePersonalCollegeApplication();
  const updateRank = useUpdatePersonalCollegeApplicationRank();
  const deleteApplication = useDeletePersonalCollegeApplication();

  const preselectedId = Number(searchParams.get("college")) || null;

  const busy =
    createApplication.isPending ||
    updateApplication.isPending ||
    updateRank.isPending ||
    deleteApplication.isPending;

  const mutationError =
    createApplication.error ??
    updateApplication.error ??
    updateRank.error ??
    deleteApplication.error;

  if (isLoading || collegesLoading) {
    return (
      <MyListShell>
        <p className="mt-6 text-sm">Loading...</p>
      </MyListShell>
    );
  }

  if (error) {
    return (
      <MyListShell>
        <p className="mt-6 text-sm">Error: {errorMessage(error)}</p>
      </MyListShell>
    );
  }

  const applications: PersonalCollegeApplication[] = applicationsData ?? [];
  const colleges: GlobalCollege[] = collegesData ?? [];

  const collegeById = new Map(colleges.map((college) => [college.id, college]));

  const collegeName = (app: PersonalCollegeApplication) =>
    collegeById.get(app.global_college_id)?.school_name ?? "Unknown college";

  // ranked applications (non null) are shown in order, unranked (null) come after them
  const ranked = applications
    .filter((application) => application.rank != null)
    .sort((a, b) => (a.rank ?? 0) - (b.rank ?? 0));

  const unranked = applications.filter(
    (application) => application.rank == null,
  );

  // colleges already on the list can't be added again
  const appliedIds = new Set(applications.map((app) => app.global_college_id));
  const availableColleges = colleges.filter(
    (college) => !appliedIds.has(college.id),
  );

  // backend answers 400 when the college has no date for the chosen deadline
  const isMissingDeadline = (
    error: ApiError,
    globalCollegeId: number,
    applicationType: string,
  ) => {
    const key = deadlineOptions.find(({ value }) => value === applicationType)
      ?.key;
    const college = collegeById.get(globalCollegeId);
    return error.status === 400 && !!key && !!college && !college[key];
  };

  // clears the last mutation's error so it doesn't linger once a new change starts
  const resetMutations = () => {
    createApplication.reset();
    updateApplication.reset();
    updateRank.reset();
    deleteApplication.reset();
  };

  // refetch full list after rank change, since the backend shifts other applications too
  const setRank = (id: number, rank: number | null) => {
    resetMutations();
    updateRank.mutate({ id, data: { rank } }, { onSuccess: () => refetch() });
  };

  // dropping on a ranked card moves the dragged application into that card's
  // position (the backend shifts everything in between). 
  // dropping on an unranked card unranks it.
  const handleDrop = (target: PersonalCollegeApplication) => {
    if (draggedId == null || draggedId === target.id) return;

    if (target.rank != null) {
      setRank(draggedId, target.rank);
      return;
    }

    const dragged = applications.find(
      (application) => application.id === draggedId,
    );
    if (dragged?.rank != null) setRank(draggedId, null);
  };

  // save category/type changes for an application, then refetch the full list to show the updated values
  const saveEdit = (
    app: PersonalCollegeApplication,
    applicationType: ApplicationType,
    category: Category,
  ) => {
    resetMutations();

    updateApplication.mutate(
      {
        id: app.id,
        data: {
          global_college_id: app.global_college_id,
          application_type: applicationType,
          category,
        },
      },
      {
        onSuccess: async () => {
          await refetch();
          setEditingId(null);
        },
        // keep the card in edit mode so the student can pick another deadline
        onError: (error) => {
          if (isMissingDeadline(error, app.global_college_id, applicationType)) {
            updateApplication.reset();
            setDeadlineError({
              globalCollegeId: app.global_college_id,
              applicationType,
            });
          }
        },
      },
    );
  };

  // add a college to the list as unranked, then refetch the full list to show it
  const handleCreate = (
    globalCollegeId: number,
    applicationType: ApplicationType,
    category: Category,
  ) => {
    resetMutations();

    createApplication.mutate(
      {
        data: {
          global_college_id: globalCollegeId,
          application_type: applicationType,
          category,
        },
      },
      {
        onSuccess: async () => {
          await refetch();
          setShowAddForm(false);
        },
        // keep the form open so the student can pick another deadline
        onError: (error) => {
          if (isMissingDeadline(error, globalCollegeId, applicationType)) {
            createApplication.reset();
            setDeadlineError({ globalCollegeId, applicationType });
          }
        },
      },
    );
  };

  // remove an application from the list, then refetch the full list to show the updated values
  const handleDelete = (app: PersonalCollegeApplication) => {
    if (!window.confirm(`Remove ${collegeName(app)} from your list?`)) return;

    resetMutations();

    deleteApplication.mutate({ id: app.id }, { onSuccess: () => refetch() });
  };

  const renderRankControls = (app: PersonalCollegeApplication) => {
    if (app.rank == null) {
      return (
        <button
          type="button"
          disabled={busy}
          onClick={() => setRank(app.id, ranked.length + 1)}
          className={smallButtonClassName}
        >
          Add to ranking
        </button>
      );
    }

    const index = ranked.indexOf(app);
    return (
      <div className="flex items-center gap-2">
        <button
          type="button"
          aria-label={`Move ${collegeName(app)} up`}
          disabled={busy || index === 0}
          onClick={() => setRank(app.id, index)}
          className={smallButtonClassName}
        >
          ↑
        </button>
        <button
          type="button"
          aria-label={`Move ${collegeName(app)} down`}
          disabled={busy || index === ranked.length - 1}
          onClick={() => setRank(app.id, index + 2)}
          className={smallButtonClassName}
        >
          ↓
        </button>
        <button
          type="button"
          disabled={busy}
          onClick={() => setRank(app.id, null)}
          className={smallButtonClassName}
        >
          Unrank
        </button>
      </div>
    );
  };

  const renderCard = (app: PersonalCollegeApplication) => {
    const isEditing = editingId === app.id;

    return (
      <ApplicationCard
        key={app.id}
        app={app}
        college={collegeById.get(app.global_college_id)}
        collegeName={collegeName(app)}
        busy={busy}
        saving={updateApplication.isPending}
        isEditing={isEditing}
        isDragging={draggedId === app.id}
        draggable={!busy && !isEditing}
        rankControls={renderRankControls(app)}
        onDragStart={() => setDraggedId(app.id)}
        onDragEnd={() => setDraggedId(null)}
        onDrop={() => handleDrop(app)}
        onStartEdit={() => {
          resetMutations();
          setEditingId(app.id);
        }}
        onCancelEdit={() => setEditingId(null)}
        onSave={(applicationType, category) =>
          saveEdit(app, applicationType, category)
        }
        onRemove={() => handleDelete(app)}
      />
    );
  };

  return (
    <MyListShell
      actions={
        <button
          type="button"
          onClick={() => setShowAddForm(!showAddForm)}
          className="text-sm underline"
        >
          {showAddForm ? "Close" : "Add a school"}
        </button>
      }
    >
      {showAddForm && (
        <AddApplicationForm
          colleges={availableColleges}
          initialCollegeId={
            // Skip the preselect once that college is already on the list
            availableColleges.some((college) => college.id === preselectedId)
              ? preselectedId
              : null
          }
          busy={busy}
          adding={createApplication.isPending}
          onAdd={handleCreate}
        />
      )}

      {mutationError && (
        <p role="alert" className="mt-3 text-sm text-red-600">
          Error: {errorMessage(mutationError)}
        </p>
      )}

      {applications.length === 0 ? (
        <p className="mt-6 text-sm">
          Your list is empty.{" "}
          <Link to="/student-college-list" className="underline">
            Browse colleges
          </Link>{" "}
          to add one.
        </p>
      ) : (
        <>
          <p className="mt-4 text-xs text-gray-500">
            Drag colleges to reorder them, or use the arrow buttons. Drop a
            college onto an unranked one to unrank it.
          </p>
          <ol className="mt-4 flex flex-col gap-4">
            {[...ranked, ...unranked].map(renderCard)}
          </ol>
        </>
      )}

      {deadlineError && (
        <DeadlineErrorDialog
          college={collegeById.get(deadlineError.globalCollegeId)}
          collegeName={
            collegeById.get(deadlineError.globalCollegeId)?.school_name ??
            "This college"
          }
          applicationType={deadlineError.applicationType}
          onClose={() => setDeadlineError(null)}
        />
      )}
    </MyListShell>
  );
}
