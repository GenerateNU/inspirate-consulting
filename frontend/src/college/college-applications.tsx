import { useState } from "react";

import { errorMessage } from "../api/client";
import {
  useListPersonalCollegeApplications,
  useListGlobalColleges,
  useUpdatePersonalCollegeApplication,
  useUpdatePersonalCollegeApplicationRank,
  useDeletePersonalCollegeApplication,
} from "../api/endpoints";
import {
  UpdatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType,
  UpdatePersonalCollegeApplicationRequestBodyCategory as Category,
  type GlobalCollege,
  type PersonalCollegeApplication,
} from "../api/models";

export default function CollegeApplications() {
  const {
    data: applicationsData,
    error,
    isLoading,
    refetch,
  } = useListPersonalCollegeApplications();
  const { data: collegesData, isLoading: collegesLoading } =
    useListGlobalColleges();

  // One hook per kind of change, the row's id is passed when calling mutate
  const updateApplication = useUpdatePersonalCollegeApplication();
  const updateRank = useUpdatePersonalCollegeApplicationRank();
  const deleteApplication = useDeletePersonalCollegeApplication();

  const busy =
    updateApplication.isPending ||
    updateRank.isPending ||
    deleteApplication.isPending;

  const mutationError =
    updateApplication.error ?? updateRank.error ?? deleteApplication.error;

  const [draggedId, setDraggedId] = useState<number | null>(null);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editType, setEditType] = useState<ApplicationType>(ApplicationType.RD);
  const [editCategory, setEditCategory] = useState<Category>(Category.target);

  if (isLoading || collegesLoading) {
    return (
      <div className="form-container">
        <h2>My Colleges</h2>
        <p>Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="form-container">
        <h2>My Colleges</h2>
        <p>Error: {errorMessage(error)}</p>
      </div>
    );
  }

  const applications: PersonalCollegeApplication[] = applicationsData ?? [];
  const colleges: GlobalCollege[] = collegesData ?? [];

  const collegeById = new Map(colleges.map((college) => [college.id, college]));

  const collegeName = (app: PersonalCollegeApplication) =>
    collegeById.get(app.global_college_id)?.school_name ?? "Unknown college";

  // Ranked applications (non null) are shown in order, unranked (null) are shown below.
  const ranked = applications
    .filter((application) => application.rank != null)
    .sort((a, b) => (a.rank ?? 0) - (b.rank ?? 0));

  const unranked = applications.filter(
    (application) => application.rank == null,
  );

  // Clears the last mutation's error so it doesn't linger once a new change starts
  const resetMutations = () => {
    updateApplication.reset();
    updateRank.reset();
    deleteApplication.reset();
  };

  // refetch full list after rank change, since the backend shifts other applications too
  const setRank = (id: number, rank: number | null) => {
    resetMutations();
    updateRank.mutate({ id, data: { rank } }, { onSuccess: () => refetch() });
  };

  // Dropping on a ranked row moves the dragged application into that row's
  // position; the backend shifts everything in between.
  const handleDropOnRow = (target: PersonalCollegeApplication) => {
    if (draggedId == null || draggedId === target.id || target.rank == null)
      return;

    setRank(draggedId, target.rank);
  };

  // Dropping on the end zone moves the application to the bottom of the
  // ranking: last place if it's already ranked, one past last if it isn't.
  const handleDropAtEnd = () => {
    if (draggedId == null) return;

    const dragged = applications.find(
      (application) => application.id === draggedId,
    );

    if (!dragged) return;

    const lastRank = dragged.rank != null ? ranked.length : ranked.length + 1;

    if (dragged.rank === lastRank) return;

    setRank(draggedId, lastRank);
  };

  const handleDropOnUnranked = () => {
    if (draggedId == null) return;

    const dragged = applications.find(
      (application) => application.id === draggedId,
    );

    if (!dragged || dragged.rank == null) return;

    setRank(draggedId, null);
  };

  // Edit an application (type and/or category) in place
  const startEditing = (app: PersonalCollegeApplication) => {
    resetMutations();

    setEditingId(app.id);

    setEditType(app.application_type as ApplicationType);
    setEditCategory(app.category as Category);
  };

  // save category/type changes for an application, then refetch the full list to show the updated values
  const saveEdit = (app: PersonalCollegeApplication) => {
    resetMutations();

    updateApplication.mutate(
      {
        id: app.id,
        data: {
          global_college_id: app.global_college_id,
          application_type: editType,
          category: editCategory,
        },
      },
      {
        onSuccess: async () => {
          await refetch();
          setEditingId(null);
        },
      },
    );
  };

  // Remove an application from the list, then refetch the full list to show the updated values
  const handleDelete = (app: PersonalCollegeApplication) => {
    if (!window.confirm(`Remove ${collegeName(app)} from your list?`)) return;

    resetMutations();

    deleteApplication.mutate({ id: app.id }, { onSuccess: () => refetch() });
  };

  const renderDetails = (app: PersonalCollegeApplication) => {
    if (editingId !== app.id) {
      return (
        <>
          <strong>{collegeName(app)}</strong> · {app.application_type} ·{" "}
          {app.category}{" "}
          <button
            type="button"
            disabled={busy}
            onClick={() => startEditing(app)}
          >
            Edit
          </button>{" "}
          <button
            type="button"
            disabled={busy}
            onClick={() => handleDelete(app)}
          >
            Remove
          </button>
        </>
      );
    }
    return (
      <>
        <strong>{collegeName(app)}</strong>{" "}
        <select
          aria-label="Application type"
          value={editType}
          onChange={(e) => setEditType(e.target.value as ApplicationType)}
        >
          {Object.values(ApplicationType).map((type) => (
            <option key={type} value={type}>
              {type}
            </option>
          ))}
        </select>{" "}
        <select
          aria-label="Category"
          value={editCategory}
          onChange={(e) => setEditCategory(e.target.value as Category)}
        >
          {Object.values(Category).map((category) => (
            <option key={category} value={category}>
              {category}
            </option>
          ))}
        </select>{" "}
        <button type="button" disabled={busy} onClick={() => saveEdit(app)}>
          {updateApplication.isPending ? "Saving..." : "Save changes"}
        </button>{" "}
        <button type="button" onClick={() => setEditingId(null)}>
          Cancel
        </button>
      </>
    );
  };

  return (
    <div className="form-container">
      <h2>My Colleges</h2>
      {mutationError && (
        <p role="alert">Error: {errorMessage(mutationError)}</p>
      )}

      <h3>Ranked</h3>
      <p>Drag colleges to reorder them, or use the arrow buttons.</p>
      {ranked.length === 0 && (
        <p>
          Nothing ranked yet. Drag a college here or choose "Add to ranking".
        </p>
      )}
      <ol>
        {ranked.map((app, index) => (
          <li
            key={app.id}
            draggable={!busy}
            onDragStart={() => setDraggedId(app.id)}
            onDragEnd={() => setDraggedId(null)}
            onDragOver={(e) => e.preventDefault()}
            onDrop={() => handleDropOnRow(app)}
            style={{ cursor: "grab", opacity: draggedId === app.id ? 0.5 : 1 }}
          >
            {renderDetails(app)}{" "}
            <button
              type="button"
              aria-label={`Move ${collegeName(app)} up`}
              disabled={busy || index === 0}
              onClick={() => setRank(app.id, index)}
            >
              ↑
            </button>
            <button
              type="button"
              aria-label={`Move ${collegeName(app)} down`}
              disabled={busy || index === ranked.length - 1}
              onClick={() => setRank(app.id, index + 2)}
            >
              ↓
            </button>{" "}
            <button
              type="button"
              disabled={busy}
              onClick={() => setRank(app.id, null)}
            >
              Unrank
            </button>
          </li>
        ))}
      </ol>
      <div
        onDragOver={(e) => e.preventDefault()}
        onDrop={handleDropAtEnd}
        style={{
          border: "1px dashed",
          padding: "0.5rem",
          marginBottom: "1rem",
        }}
      >
        Drop here to move to the bottom of your ranking
      </div>

      <h3>Not ranked</h3>
      <div
        onDragOver={(e) => e.preventDefault()}
        onDrop={handleDropOnUnranked}
        style={{ border: "1px dashed", padding: "0.5rem" }}
      >
        {unranked.length === 0 && (
          <p>
            Every college on your list is ranked. Drop one here to unrank it.
          </p>
        )}
        <ul>
          {unranked.map((app) => (
            <li
              key={app.id}
              draggable={!busy}
              onDragStart={() => setDraggedId(app.id)}
              onDragEnd={() => setDraggedId(null)}
              style={{
                cursor: "grab",
                opacity: draggedId === app.id ? 0.5 : 1,
              }}
            >
              {renderDetails(app)}{" "}
              <button
                type="button"
                disabled={busy}
                onClick={() => setRank(app.id, ranked.length + 1)}
              >
                Add to ranking
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
