import { useState } from "react";
import { errorMessage } from "../api/client";
import {
  CATEGORY_LABELS,
  STATUS_LABELS,
  type ExtracurricularFormValues,
} from "./ExtracurricularShared";

type Props = {
  initialValues: ExtracurricularFormValues;
  submitLabel: string;
  pendingLabel: string;
  isPending: boolean;
  error: Error | null;
  onSubmit: (values: ExtracurricularFormValues) => void;
};

const endsBeforeStart = (values: ExtracurricularFormValues) =>
  values.end_date !== "" &&
  values.start_date !== "" &&
  values.end_date < values.start_date;

export default function ExtracurricularForm({
  initialValues,
  submitLabel,
  pendingLabel,
  isPending,
  error,
  onSubmit,
}: Props) {
  const [values, setValues] = useState(initialValues);
  const [dateError, setDateError] = useState<string | null>(null);

  const update = <K extends keyof ExtracurricularFormValues>(
    key: K,
    value: ExtracurricularFormValues[K],
  ) => {
    setValues((prev) => ({ ...prev, [key]: value }));
    if (key === "start_date" || key === "end_date") setDateError(null);
  };

  const handleSubmit = (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (endsBeforeStart(values)) {
      setDateError("End date can't be before the start date.");
      return;
    }
    onSubmit(values);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="ec-section">
        <h3>Basics</h3>
        <div className="ec-form-grid">
          <label className="ec-field">
            Activity name
            <input
              type="text"
              placeholder="e.g. Track & Field"
              value={values.name}
              onChange={(e) => update("name", e.target.value)}
              required
            />
          </label>
          <label className="ec-field">
            Position / role (optional)
            <input
              type="text"
              placeholder="Max 50 characters"
              maxLength={50}
              value={values.leadership_role}
              onChange={(e) => update("leadership_role", e.target.value)}
            />
          </label>
          <label className="ec-field">
            Organization
            <input
              type="text"
              placeholder="Max 100 characters"
              maxLength={100}
              value={values.organization}
              onChange={(e) => update("organization", e.target.value)}
              required
            />
          </label>
          <label className="ec-field">
            Category
            <select
              value={values.type}
              onChange={(e) =>
                update(
                  "type",
                  e.target.value as ExtracurricularFormValues["type"],
                )
              }
            >
              {Object.entries(CATEGORY_LABELS).map(([value, label]) => (
                <option key={value} value={value}>
                  {label}
                </option>
              ))}
            </select>
          </label>
        </div>
      </div>

      <div className="ec-section">
        <h3>Timeline</h3>
        <div className="ec-form-grid">
          <label className="ec-field">
            Status
            <select
              value={values.status}
              onChange={(e) =>
                update(
                  "status",
                  e.target.value as ExtracurricularFormValues["status"],
                )
              }
            >
              {Object.entries(STATUS_LABELS).map(([value, label]) => (
                <option key={value} value={value}>
                  {label}
                </option>
              ))}
            </select>
          </label>
          <label className="ec-field">
            Start date
            <input
              type="date"
              value={values.start_date}
              onChange={(e) => update("start_date", e.target.value)}
              required
            />
          </label>
          <label className="ec-field">
            End date (optional)
            <input
              type="date"
              min={values.start_date || undefined}
              value={values.end_date}
              onChange={(e) => update("end_date", e.target.value)}
              aria-invalid={dateError ? true : undefined}
              aria-describedby={dateError ? "ec-date-error" : undefined}
            />
            {dateError && (
              <span id="ec-date-error" role="alert">
                {dateError}
              </span>
            )}
          </label>
        </div>
      </div>

      <div className="ec-section">
        <h3>Application description</h3>
        <label className="ec-field ec-field-wide">
          Up to 250 words (Common App)
          <textarea
            rows={6}
            placeholder="Write your description here..."
            value={values.description}
            onChange={(e) => update("description", e.target.value)}
            required
          />
        </label>
      </div>

      {error && <p>Error: {errorMessage(error)}</p>}
      <div className="ec-actions">
        <button type="submit" className="ec-button" disabled={isPending}>
          {isPending ? pendingLabel : submitLabel}
        </button>
      </div>
    </form>
  );
}