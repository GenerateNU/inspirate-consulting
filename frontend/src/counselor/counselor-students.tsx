import { useState } from "react";

import { errorMessage } from "../api/client";
import { useListCounselorStudents } from "../api/endpoints";
import {
  ListCounselorStudentsYear as Year,
  ListCounselorStudentsSortBy as SortBy,
  ListCounselorStudentsSortOrder as SortOrder,
  type CounselorStudent,
} from "../api/models";

export default function CounselorStudents() {
  
  // Filters and sorting are sent as query params, so changing them refetches the list
  const [year, setYear] = useState<Year | "">("");
  const [pendingReviewOnly, setPendingReviewOnly] = useState(false);
  const [sortBy, setSortBy] = useState<SortBy>(SortBy.name);
  const [sortOrder, setSortOrder] = useState<SortOrder>(SortOrder.asc);

  const {
    data: studentsData,
    error,
    isLoading,
  } = useListCounselorStudents({
    year: year || undefined,
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
          Year:{" "}
          <select
            value={year}
            onChange={(e) => setYear(e.target.value as typeof year)}
          >
            <option value="">All years</option>
            {Object.values(Year).map((y) => (
              <option key={y} value={y}>
                {y}
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
                {field}
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
            · {student.year}
            {student.has_pending_review && " · Review requested"}
            {student.is_falling_behind && " · Falling behind"}
          </li>
        ))}
      </ul>
    </div>
  );
}
