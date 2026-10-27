import { useEffect, useRef } from "react";

import type { GlobalCollege } from "../../api/models";
import { deadlineOptions, formatDate, primaryButtonClassName } from "../shared";
import { typeLabel } from "./application-labels";

// popup shown when the backend rejects an application because the college doesn't offer the chosen deadline
// lists the deadlines it does offer
export default function DeadlineErrorDialog({
  college,
  collegeName,
  applicationType,
  onClose,
}: {
  college: GlobalCollege | undefined;
  collegeName: string;
  applicationType: string;
  onClose: () => void;
}) {
  const dialogRef = useRef<HTMLDialogElement>(null);


  useEffect(() => {
    const dialog = dialogRef.current;
    if (dialog && !dialog.open) dialog.showModal();
  }, []);

  const offered = college
    ? deadlineOptions.filter(({ key }) => college[key])
    : [];

  return (
    <dialog
      ref={dialogRef}
      onClose={onClose}
      aria-labelledby="deadline-error-title"
      className="m-auto w-full max-w-md border border-black p-6 backdrop:bg-black/40"
    >
      <h3 id="deadline-error-title" className="mb-2 text-xl font-bold">
        Deadline not available
      </h3>
      <p className="text-sm">
        {collegeName} doesn&apos;t offer {typeLabel(applicationType)} (
        {applicationType}).
      </p>

      {offered.length > 0 ? (
        <>
          <p className="mt-4 text-sm">Pick one of the deadlines it offers:</p>
          <ul className="mt-2 list-disc pl-5 text-sm">
            {offered.map(({ value, key, label }) => (
              <li key={value}>
                {label} ({value}): {formatDate(college![key]!)}
              </li>
            ))}
          </ul>
        </>
      ) : (
        <p className="mt-4 text-sm">
          This college doesn&apos;t have any deadlines set yet. Ask your
          counselor to add them.
        </p>
      )}

      <div className="mt-6 flex justify-end">
        <button type="button" onClick={onClose} className={primaryButtonClassName}>
          OK
        </button>
      </div>
    </dialog>
  );
}
