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

// PLACEHOLDER: replace with the counselor's student list once there's an endpoint/hook for it
const PLACEHOLDER_STUDENTS = [
  { id: "placeholder-1", name: "Risha Jones" },
  { id: "placeholder-2", name: "Hasith Klein" },
  { id: "placeholder-3", name: "Shrikar Rodrigo" },
  { id: "placeholder-4", name: "Abhinav Smith" },
];

function toggle(list: string[], item: string) {
  return list.includes(item) ? list.filter((x) => x !== item) : [...list, item];
}

export default function AssignToSelect({
  value,
  onChange,
}: AssignToSelectProps) {
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
          <p className="text-xs text-gray-500">Placeholder student list</p>
          {PLACEHOLDER_STUDENTS.map((student) => (
            <label
              key={student.id}
              className="flex items-center justify-between text-sm"
            >
              {student.name}
              <input
                type="checkbox"
                checked={value.studentIds.includes(student.id)}
                onChange={() =>
                  onChange({
                    ...value,
                    studentIds: toggle(value.studentIds, student.id),
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
