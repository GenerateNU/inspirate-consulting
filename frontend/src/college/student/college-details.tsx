import { Link, useParams } from "react-router-dom";

import { errorMessage, type ApiError } from "../../api/client";
import { useGetGlobalCollege } from "../../api/endpoints";
import { Field, deadlineOptions, formatDate } from "../shared";

// read-only view of one college
export default function StudentCollegeDetails() {
  const id = Number(useParams().id);
  const { data: college, isLoading, error } = useGetGlobalCollege(id);

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
        <dl>
          <Field label="Location">{college.school_location}</Field>
          {deadlineOptions.map(({ value, key, label }) => (
            <Field key={value} label={`${label} (${value})`}>
              {college[key] ? formatDate(college[key]!) : "N/A"}
            </Field>
          ))}
        </dl>
      </>
    );
  }

  return (
    <div className="mx-auto w-full max-w-5xl px-4 py-10">
      <Link
        to="/student-college-list"
        className="mb-4 inline-block text-sm hover:underline"
      >
        ← Back to colleges
      </Link>
      <div className="border border-black p-6">{body}</div>
    </div>
  );
}
