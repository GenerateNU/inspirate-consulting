import { errorMessage } from "../../api/client";
import CollegeDetailsLayout, {
  useCollegeFromRoute,
} from "../utils/college-details-layout";
import Field from "../utils/components/field";
import {
  deadlineOptions,
  inputClassName,
  primaryButtonClassName,
} from "../utils/shared";
import {
  useCollegeEditForm,
  type CollegeEditFormProps,
} from "./use-college-edit-form";

// counselor-only form for updating a college
function CollegeEditForm(props: CollegeEditFormProps) {
  const form = useCollegeEditForm(props);

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        form.save();
      }}
    >
      <dl>
        <Field label="School name">
          <input
            type="text"
            value={form.name}
            onChange={(e) => form.setName(e.target.value)}
            required
            className={inputClassName}
          />
        </Field>
        <Field label="Location">
          <input
            type="text"
            value={form.location}
            onChange={(e) => form.setLocation(e.target.value)}
            required
            className={inputClassName}
          />
        </Field>
        {deadlineOptions.map(({ value, key, label }) => (
          <Field key={value} label={`${label} (${value})`}>
            <input
              type="date"
              value={form.deadlines[key]}
              onChange={(e) => form.setDeadline(key, e.target.value)}
              className={inputClassName}
            />
          </Field>
        ))}
      </dl>

      <div className="mt-6 flex flex-wrap items-center gap-3">
        <button
          type="submit"
          disabled={form.busy}
          className={primaryButtonClassName}
        >
          {form.saving ? "Saving..." : "Save changes"}
        </button>
        {form.isDuplicateWarning && (
          <button
            type="button"
            disabled={form.busy}
            onClick={() => form.save(true)}
            className="rounded-md border border-black px-5 py-2 text-sm hover:bg-gray-100"
          >
            Save anyway
          </button>
        )}
        <button
          type="button"
          disabled={form.busy}
          onClick={form.deleteCollege}
          className="ml-auto rounded-md border border-red-600 px-5 py-2 text-sm text-red-600 hover:bg-red-50 disabled:opacity-40"
        >
          {form.deleting ? "Deleting..." : "Delete college"}
        </button>
      </div>
      {form.status && <p className="mt-3 text-sm">{form.status}</p>}
      {/* The backend's 409 message is developer-facing, so show our own wording */}
      {form.isDuplicateWarning ? (
        <p
          role="alert"
          className="mt-3 border border-amber-300 bg-amber-100 p-3 text-sm text-amber-700"
        >
          &quot;{form.name}&quot; looks very similar to another college on the
          list. Check that it isn&apos;t a duplicate, then choose Save anyway to
          save your changes.
        </p>
      ) : (
        form.mutationError && (
          <p role="alert" className="mt-3 text-sm text-red-600">
            Error: {errorMessage(form.mutationError)}
          </p>
        )
      )}
    </form>
  );
}

// loads the college, CollegeEditForm needs loaded data so its fields can start from the saved values
export default function CounselorCollegeDetails() {
  const { data: college, isLoading, error, refetch } = useCollegeFromRoute();

  return (
    <CollegeDetailsLayout college={college} isLoading={isLoading} error={error}>
      {college && <CollegeEditForm college={college} onSaved={refetch} />}
    </CollegeDetailsLayout>
  );
}
