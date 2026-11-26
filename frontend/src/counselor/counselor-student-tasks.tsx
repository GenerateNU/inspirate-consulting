import { Link, useParams } from "react-router-dom";
import { useGetStudent, useFetchUser } from "../api/endpoints";
import {
  GetTodoItemsSortBy,
  GetTodoItemsSortOrder,
  GetTodoItemsStatus,
} from "../api/models";
import { graduationYearLabel } from "../lib/graduation-year";
import { formatDeadline } from "../task/todoShared";
import { useTodoPages } from "../task/useTodoList";
import { ErrorMessage } from "@/components/ui/error-message";

//validate URL before anything is fetched
export default function CounselorStudentTasks() {
  const { id } = useParams();

  if (!id) {
    return (
      <div className="p-10">
        <ErrorMessage
          variant="banner"
          message="No student was specified."
          action={
            <Link to="/counselor-students" className="underline">
              Back to students
            </Link>
          }
        />
      </div>
    );
  }
  return <CounselorStudentTasksContent studentId={id} />;
}

//only rendered with studentID
function CounselorStudentTasksContent({ studentId }: { studentId: string }) {
  const studentQuery = useGetStudent(studentId, {
    query: { enabled: studentId !== "" },
  });
  const student = studentQuery.data;

  const userId = student?.user_id ?? "";
  const userQuery = useFetchUser(userId, { query: { enabled: userId !== "" } });
  const user = userQuery.data;

  if (studentQuery.isLoading || userQuery.isLoading) return <p>Loading...</p>;
  if (studentQuery.error || !student)
    return (
      <div className="p-10">
        <ErrorMessage
          variant="banner"
          title="Couldn't load this student"
          error={studentQuery.error}
          message={studentQuery.error ? undefined : "Student not found."}
          onRetry={studentQuery.error ? studentQuery.refetch : undefined}
        />
      </div>
    );

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
            <li>
              Graduation Year: {graduationYearLabel(student.graduation_year)}
            </li>
            <li>High School: {student.organization ?? "Not set"}</li>
            <li>Major of Interest: {student.student_interest ?? "Not set"}</li>
            <li>Review Balance: {student.review_balance}</li>
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
              to={`/assign-task?studentId=${studentId}`}
              className="rounded bg-black px-6 py-2 text-white"
            >
              Assign Task →
            </Link>
          </div>

          <StudentTaskRows studentId={studentId} />
        </div>
      </div>
    </div>
  );
}

// This student's tasks, newest first, each with a link to edit it
function StudentTaskRows({ studentId }: { studentId: string }) {
  const { items, error, isLoading, hasMore, isLoadingMore, loadMore, refresh } =
    useTodoPages(
      {
        student_id: studentId,
        status: GetTodoItemsStatus.all,
        sort_by: GetTodoItemsSortBy.created_at,
        sort_order: GetTodoItemsSortOrder.desc,
      },
      "",
    );

  if (isLoading) return <p className="text-sm">Loading tasks...</p>;
  if (error)
    return (
      <ErrorMessage
        variant="banner"
        title="Couldn't load tasks"
        error={error}
        onRetry={refresh}
      />
    );
  if (items.length === 0)
    return <p className="text-sm text-gray-500">No tasks assigned yet.</p>;

  return (
    <div className="flex flex-col gap-2">
      <ul className="flex flex-col gap-2">
        {items.map((task) => (
          <li
            key={task.id}
            className="flex items-center justify-between gap-4 rounded-lg bg-white px-4 py-3"
          >
            <div className="min-w-0">
              <p
                className={`truncate text-sm font-medium ${task.completed_at ? "text-gray-500 line-through" : ""}`}
              >
                {task.title}
              </p>
              {task.deadline && (
                <p className="text-xs text-gray-500">
                  Due {formatDeadline(task.deadline)}
                </p>
              )}
            </div>
            <Link
              to={`/counselor/students/${studentId}/tasks/${task.id}/edit`}
              className="shrink-0 text-sm underline"
            >
              Edit
            </Link>
          </li>
        ))}
      </ul>
      {hasMore && (
        <button
          type="button"
          onClick={() => void loadMore()}
          disabled={isLoadingMore}
          className="self-start px-4 py-1 text-sm"
        >
          {isLoadingMore ? "Loading..." : "Load more"}
        </button>
      )}
    </div>
  );
}
