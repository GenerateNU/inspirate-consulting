import { useState } from "react";
import { Link } from "react-router-dom";

import { errorMessage } from "../api/client";
import {
  graduationYearLabel,
  graduationYearOptions,
} from "@/lib/graduation-year";
import { useListCounselorStudents } from "../api/endpoints";
import {
  ListCounselorStudentsSortBy as SortBy,
  ListCounselorStudentsSortOrder as SortOrder,
  type CounselorStudent,
} from "../api/models";

export default function CounselorStudents() {
  // Filters and sorting are sent as query params, so changing them refetches the list
  const [graduationYear, setGraduationYear] = useState<number | "">("");
  const [pendingReviewOnly, setPendingReviewOnly] = useState(false);
  const [sortBy, setSortBy] = useState<SortBy>(SortBy.name);
  const [sortOrder, setSortOrder] = useState<SortOrder>(SortOrder.asc);

  const {
    data: studentsData,
    error,
    isLoading,
  } = useListCounselorStudents({
    graduation_year: graduationYear || undefined,
    has_pending_review: pendingReviewOnly || undefined,
    sort_by: sortBy,
    sort_order: sortOrder,
  });

  const students: CounselorStudent[] = studentsData ?? [];

  return (
    <div className="form-container">
      <h2>My Students</h2>

      <div className="form-group">
        <label>
          Graduation year:{" "}
          <select
            value={graduationYear}
            onChange={(e) =>
              setGraduationYear(
                e.target.value === "" ? "" : Number(e.target.value),
              )
            }
          >
            <option value="">All years</option>
            {graduationYearOptions().map(({ value, label }) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
        </label>{" "}
        <label>
          <input
            type="checkbox"
            checked={pendingReviewOnly}
            onChange={(e) => setPendingReviewOnly(e.target.checked)}
          />{" "}
          Only students with a pending review
        </label>
      </div>

      <div className="form-group">
        <label>
          Sort by:{" "}
          <select
            value={sortBy}
            onChange={(e) => setSortBy(e.target.value as SortBy)}
          >
            {Object.values(SortBy).map((field) => (
              <option key={field} value={field}>
                {field.replace("_", " ")}
              </option>
            ))}
          </select>
        </label>{" "}
        <select
          aria-label="Sort order"
          value={sortOrder}
          onChange={(e) => setSortOrder(e.target.value as SortOrder)}
        >
          <option value={SortOrder.asc}>Ascending</option>
          <option value={SortOrder.desc}>Descending</option>
        </select>
      </div>

      {isLoading && <p>Loading...</p>}
      {error && <p>Error: {errorMessage(error)}</p>}
      {!isLoading && !error && students.length === 0 && (
        <p>No students match these filters.</p>
      )}

      <ul>
        {students.map((student) => (
          <li key={student.student_id}>
            {/* Red marks a requested essay review, placeholder until the design is final */}
            <strong
              style={{ color: student.has_pending_review ? "red" : undefined }}
            >
              {student.name}
            </strong>{" "}
            · {graduationYearLabel(student.graduation_year)}
            {student.has_pending_review && " · Review requested"}
            {student.is_falling_behind && " · Falling behind"}
            {" · "}
            <Link to={`/counselor/students/${student.student_id}/essays`}>
              Essays
            </Link>
            {" · "}
            <Link to={`/counselor/students/${student.student_id}/colleges`}>
              Colleges
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
