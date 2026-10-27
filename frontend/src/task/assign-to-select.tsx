import { useListCounselorStudents } from "../api/endpoints";
import type { CounselorStudent } from "../api/models";

export type AssignToMode = "" | "all" | "grade" | "individual";

export type AssignToValue = {
  mode: AssignToMode;
  grades: string[];
  studentIds: string[];
};

type AssignToSelectProps = {
  value: AssignToValue;
  onChange: (value: AssignToValue) => void;
};

const GRADES = ["9th", "10th", "11th", "12th"];

function toggle(list: string[], item: string) {
  return list.includes(item) ? list.filter((x) => x !== item) : [...list, item];
}

export default function AssignToSelect({
  value,
  onChange,
}: AssignToSelectProps) {
  // Fetch the counselor's students, only when "Individual" is selected
  const studentsQuery = useListCounselorStudents(
    {},
    { query: { enabled: value.mode === "individual" } },
  );
  const students: CounselorStudent[] = studentsQuery.data ?? [];

  return (
    <div className="flex flex-col gap-2">
      <label className="text-sm">
        Assign to
        <select
          value={value.mode}
          onChange={(e) =>
            onChange({ ...value, mode: e.target.value as AssignToMode })
          }
          className="mt-1 block w-full border px-2 py-1"
        >
          <option value="">Select an option</option>
          <option value="all">All</option>
          <option value="grade">Grade Level</option>
          <option value="individual">Individual</option>
        </select>
      </label>

      {value.mode === "grade" && (
        <div className="flex flex-col gap-1 border p-2">
          {GRADES.map((grade) => (
            <label
              key={grade}
              className="flex items-center justify-between text-sm"
            >
              {grade}
              <input
                type="checkbox"
                checked={value.grades.includes(grade)}
                onChange={() =>
                  onChange({ ...value, grades: toggle(value.grades, grade) })
                }
              />
            </label>
          ))}
        </div>
      )}

      {value.mode === "individual" && (
        <div className="flex flex-col gap-1 border p-2">
          {studentsQuery.isLoading && (
            <p className="text-sm">Loading students...</p>
          )}
          {studentsQuery.error && (
            <p className="text-sm">Could not load students.</p>
          )}
          {!studentsQuery.isLoading &&
            !studentsQuery.error &&
            students.length === 0 && (
              <p className="text-sm">No students found.</p>
            )}

          {students.map((student) => (
            <label
              key={student.student_id}
              className="flex items-center justify-between text-sm"
            >
              {student.name}
              <input
                type="checkbox"
                checked={value.studentIds.includes(student.student_id)}
                onChange={() =>
                  onChange({
                    ...value,
                    studentIds: toggle(value.studentIds, student.student_id),
                  })
                }
              />
            </label>
          ))}
        </div>
      )}
    </div>
  );
}