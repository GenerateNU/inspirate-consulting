import { useEffect, useRef, useState } from "react";

import { errorMessage, type ApiError } from "../../api/client";
import { useCreateGlobalCollege } from "../../api/endpoints";
import Field from "../utils/components/field";
import { inputClassName, primaryButtonClassName } from "../utils/shared";
import CyclesEditor, {
  cyclesDraftError,
  cyclesFromDraft,
  draftFromCycles,
  type CyclesDraft,
} from "./cycles-editor";

// popup form for adding a global college
interface AddCollegeModalProps {
  onClose: () => void;
  onCreated: () => Promise<unknown>;
}

export default function AddCollegeModal({
  onClose,
  onCreated,
}: AddCollegeModalProps) {
  const [name, setName] = useState("");
  const [location, setLocation] = useState("");
  const [cycles, setCycles] = useState<CyclesDraft>(() =>
    draftFromCycles(null),
  );
  const [validationError, setValidationError] = useState<string | null>(null);

  const dialogRef = useRef<HTMLDialogElement>(null);

  const createCollege = useCreateGlobalCollege();

  useEffect(() => {
    const dialog = dialogRef.current;
    if (dialog && !dialog.open) dialog.showModal();
  }, []);

  // 409 means the new college is a near-duplicate of an existing one
  const isDuplicateWarning =
    (createCollege.error as ApiError | null)?.status === 409;

  const create = (confirmDuplicate = false) => {
    createCollege.reset();
    const problem = cyclesDraftError(cycles);
    setValidationError(problem);
    if (problem) return;

    createCollege.mutate(
      {
        data: {
          school_name: name,
          school_location: location,
          cycles: cyclesFromDraft(cycles),
          confirm_duplicate: confirmDuplicate,
        },
      },
      {
        onSuccess: async () => {
          await onCreated();
          onClose();
        },
      },
    );
  };

  return (
    <dialog
      ref={dialogRef}
      onClose={onClose}
      className="m-auto w-full max-w-3xl border border-black p-6 backdrop:bg-black/40"
    >
      <form
        onSubmit={(e) => {
          e.preventDefault();
          create();
        }}
      >
        <h3 className="mb-4 text-2xl font-bold">Add a college</h3>

        <dl>
          <Field label="School name *">
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              className={inputClassName}
            />
          </Field>
          <Field label="Location *">
            <input
              type="text"
              value={location}
              onChange={(e) => setLocation(e.target.value)}
              required
              className={inputClassName}
            />
          </Field>
        </dl>

        <h4 className="mt-6 mb-2 font-host text-sm font-medium">
          Application cycles &amp; deadlines
        </h4>
        <CyclesEditor
          draft={cycles}
          disabled={createCollege.isPending}
          onChange={setCycles}
        />

        {/* The backend's 409 message is developer-facing, so show our own wording */}
        {validationError ? (
          <p role="alert" className="mt-3 text-sm text-red-600">
            {validationError}
          </p>
        ) : isDuplicateWarning ? (
          <p
            role="alert"
            className="mt-3 border border-amber-300 bg-amber-100 p-3 text-sm text-amber-700"
          >
            &quot;{name}&quot; looks very similar to a college that&apos;s
            already on the list. Check that it isn&apos;t a duplicate, then
            choose Add anyway to add it.
          </p>
        ) : (
          createCollege.error && (
            <p role="alert" className="mt-3 text-sm text-red-600">
              Error: {errorMessage(createCollege.error)}
            </p>
          )
        )}

        <div className="mt-6 flex flex-wrap items-center justify-end gap-3">
          <button
            type="button"
            onClick={onClose}
            className="rounded-md border border-black px-5 py-2 text-sm hover:bg-gray-100"
          >
            Cancel
          </button>
          {isDuplicateWarning && (
            <button
              type="button"
              disabled={createCollege.isPending}
              onClick={() => create(true)}
              className="rounded-md border border-black px-5 py-2 text-sm hover:bg-gray-100"
            >
              Add anyway
            </button>
          )}
          <button
            type="submit"
            disabled={createCollege.isPending}
            className={primaryButtonClassName}
          >
            {createCollege.isPending ? "Adding..." : "Add College"}
          </button>
        </div>
      </form>
    </dialog>
  );
}
