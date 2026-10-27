import { useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import { errorMessage, type ApiError } from "../../api/client";
import {
  useDeleteGlobalCollege,
  useGetGlobalCollege,
  useUpdateGlobalCollege,
} from "../../api/endpoints";
import type { GlobalCollege } from "../../api/models";
import {
  Field,
  deadlineOptions,
  inputClassName,
  primaryButtonClassName,
} from "../shared";

const backTo = "/counselor-college-list";

const toISODateTime = (dateStr: string): string | undefined =>
  dateStr ? `${dateStr}T00:00:00Z` : undefined;

// date inputs need YYYY-MM-DD
const toDateInput = (dateTime?: string | null): string =>
  dateTime ? dateTime.slice(0, 10) : "";

// counselor-only form for updating a college
function CollegeEditForm({
  college,
  onSaved,
}: {
  college: GlobalCollege;
  onSaved: () => Promise<unknown>;
}) {
  const [name, setName] = useState(college.school_name);
  const [location, setLocation] = useState(college.school_location);
  const [deadlines, setDeadlines] = useState({
    ea_deadline: toDateInput(college.ea_deadline),
    ed_deadline: toDateInput(college.ed_deadline),
    rd_deadline: toDateInput(college.rd_deadline),
  });
  const [status, setStatus] = useState("");

  const navigate = useNavigate();

  const updateCollege = useUpdateGlobalCollege();
  const deleteCollege = useDeleteGlobalCollege();

  const busy = updateCollege.isPending || deleteCollege.isPending;
  const mutationError = updateCollege.error ?? deleteCollege.error;

  // 409 means the edit is a near-duplicate of another college
  const isDuplicateWarning = (updateCollege.error as ApiError | null)?.status === 409;

  // clears the last mutation's error
  const resetMutations = () => {
    setStatus("");
    updateCollege.reset();
    deleteCollege.reset();
  };

  // save changes, then refetch
  const save = (confirmDuplicate = false) => {
    resetMutations();

    updateCollege.mutate(
      {
        id: college.id,
        data: {
          school_name: name,
          school_location: location,
          ea_deadline: toISODateTime(deadlines.ea_deadline),
          ed_deadline: toISODateTime(deadlines.ed_deadline),
          rd_deadline: toISODateTime(deadlines.rd_deadline),
          confirm_duplicate: confirmDuplicate,
        },
      },
      {
        onSuccess: async () => {
          await onSaved();
          setStatus("Saved.");
        },
      },
    );
  };

  // delete the college, then go back to the list since this page is gone
  const handleDelete = () => {
    if (!window.confirm(`Delete ${college.school_name}?`)) return;

    resetMutations();

    deleteCollege.mutate(
      { id: college.id },
      { onSuccess: () => navigate(backTo) },
    );
  };

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        save();
      }}
    >
      <dl>
        <Field label="School name">
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            className={inputClassName}
          />
        </Field>
        <Field label="Location">
          <input
            type="text"
            value={location}
            onChange={(e) => setLocation(e.target.value)}
            required
            className={inputClassName}
          />
        </Field>
        {deadlineOptions.map(({ value, key, label }) => (
          <Field key={value} label={`${label} (${value})`}>
            <input
              type="date"
              value={deadlines[key]}
              onChange={(e) =>
                setDeadlines({ ...deadlines, [key]: e.target.value })
              }
              className={inputClassName}
            />
          </Field>
        ))}
      </dl>

      <div className="mt-6 flex flex-wrap items-center gap-3">
        <button type="submit" disabled={busy} className={primaryButtonClassName}>
          {updateCollege.isPending ? "Saving..." : "Save changes"}
        </button>
        {isDuplicateWarning && (
          <button
            type="button"
            disabled={busy}
            onClick={() => save(true)}
            className="rounded-md border border-black px-5 py-2 text-sm hover:bg-gray-100"
          >
            Save anyway
          </button>
        )}
        <button
          type="button"
          disabled={busy}
          onClick={handleDelete}
          className="ml-auto rounded-md border border-red-600 px-5 py-2 text-sm text-red-600 hover:bg-red-50 disabled:opacity-40"
        >
          {deleteCollege.isPending ? "Deleting..." : "Delete college"}
        </button>
      </div>
      {status && <p className="mt-3 text-sm">{status}</p>}
      {/* The backend's 409 message is developer-facing, so show our own wording */}
      {isDuplicateWarning ? (
        <p
          role="alert"
          className="mt-3 border border-amber-300 bg-amber-100 p-3 text-sm text-amber-700"
        >
          &quot;{name}&quot; looks very similar to another college on the list.
          Check that it isn&apos;t a duplicate, then choose Save anyway to save
          your changes.
        </p>
      ) : (
        mutationError && (
          <p role="alert" className="mt-3 text-sm text-red-600">
            Error: {errorMessage(mutationError)}
          </p>
        )
      )}
    </form>
  );
}

// loads the college, CollegeEditForm needs loaded data so its fields can start from the saved values
export default function CounselorCollegeDetails() {
  const id = Number(useParams().id);
  const { data: college, isLoading, error, refetch } = useGetGlobalCollege(id);

  let body;
  if (isLoading) {
    body = <p className="text-sm">Loading...</p>;
  } else if (error || !college) {
    body = (
      <p className="text-sm">
        {(error as ApiError | null)?.status === 404
          ? "College not found."
          : `Error: ${errorMessage(error)}`}
      </p>
    );
  } else {
    body = (
      <>
        <h3 className="mb-4 text-3xl font-bold">{college.school_name}</h3>
        <CollegeEditForm college={college} onSaved={refetch} />
      </>
    );
  }

  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-10">
      <Link to={backTo} className="mb-4 inline-block text-sm hover:underline">
        ← Back to colleges
      </Link>
      <div className="border border-black p-6">{body}</div>
    </div>
  );
}
