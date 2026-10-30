import { Link, useParams } from "react-router-dom";
import { useGetStudent, useFetchUser } from "../api/endpoints";

const gradeByYear: Record<string, number> = {
  freshman: 9,
  sophomore: 10,
  junior: 11,
  senior: 12,
};

export default function CounselorStudentTasks() {
  const { id = "" } = useParams();

  const studentQuery = useGetStudent(id, { query: { enabled: id !== "" } });
  const student = studentQuery.data;

  const userId = student?.user_id ?? "";
  const userQuery = useFetchUser(userId, { query: { enabled: userId !== "" } });
  const user = userQuery.data;

  if (studentQuery.isLoading || userQuery.isLoading) return <p>Loading...</p>;
  if (studentQuery.error || !student) return <p>Could not load student.</p>;

  const name = user?.name ?? "Student";

  return (
    <div className="flex flex-col gap-6 p-10">
      <h1 className="text-2xl font-semibold">{name}&apos;s Tasks</h1>

      <div className="flex gap-8">
        {/* Left: profile card */}
        <div className="flex w-56 flex-col items-center gap-2 bg-white p-4">
          {/* PLACEHOLDER: profile photo */}
          <div className="h-20 w-20 rounded-full bg-gray-300" />
          <p className="font-semibold">{name}</p>

          <p className="mt-2 text-sm font-semibold">Student Info</p>
          <ul className="text-xs">
            <li>Grade: {gradeByYear[student.year] ?? student.year}</li>
            <li>High School: {student.organization ?? "Not set"}</li>
            <li>Major of Interest: {student.student_interest ?? "Not set"}</li>
          </ul>

          {/* PLACEHOLDER: progress status needs this student's tasks */}
          <p className="mt-2 text-sm font-semibold">Progress Status</p>
          <p className="text-xs text-gray-500">Not available yet</p>
        </div>

        {/* Right: legend, assign button, task list */}
        <div className="flex flex-1 flex-col gap-4">
          <div className="flex items-center justify-between">
            <div className="flex gap-12 text-xs">
              <span>Extracurriculars</span>
              <span>Essays</span>
              <span>Modules</span>
            </div>

            <Link
              to={`/assign-task?studentId=${id}`}
              className="rounded bg-black px-6 py-2 text-white"
            >
              Assign Task →
            </Link>
          </div>

          {/* PLACEHOLDER: needs an endpoint for a specific student's tasks */}
          <div className="border p-6 text-sm text-gray-500">
            This student&apos;s task list isn&apos;t available yet. It needs a
            backend endpoint for listing a specific student&apos;s tasks.
          </div>
        </div>
      </div>
    </div>
  );
}
