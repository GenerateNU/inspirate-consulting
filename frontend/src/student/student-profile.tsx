// name, grade, highschool, major of interest, tasks completed
import { useParams } from "react-router-dom";
import { useGetStudent, useFetchUser, useGetTodoItems } from "../api/endpoints";
import { graduationYearLabel } from "../lib/graduation-year";
import { ErrorMessage } from "@/components/ui/error-message";

export default function StudentProfile() {
  const { id = "" } = useParams();

  // Fetch the student using the ID from the URL
  const studentQuery = useGetStudent(id, {
    query: { enabled: id !== "" },
  });
  const student = studentQuery.data;

  // Fetch the user once the student's user_id is known
  const userId = student?.user_id ?? "";
  const userQuery = useFetchUser(userId, {
    query: { enabled: userId !== "" },
  });
  const user = userQuery.data;

  //Fetch todo items to count completed tasks
  const todosQuery = useGetTodoItems();
  const todos = Array.isArray(todosQuery.data) ? todosQuery.data : [];

  if (studentQuery.isLoading || userQuery.isLoading || todosQuery.isLoading) {
    return <p>Loading...</p>;
  }

  const failed = [
    { query: studentQuery, title: "Couldn't load student" },
    { query: userQuery, title: "Couldn't load user" },
    { query: todosQuery, title: "Couldn't load tasks" },
  ].find(({ query }) => query.error);
  if (failed) {
    return (
      <ErrorMessage
        variant="banner"
        title={failed.title}
        error={failed.query.error}
        onRetry={failed.query.refetch}
        className="m-8 w-auto"
      />
    );
  }

  if (!student) {
    return (
      <ErrorMessage
        variant="banner"
        message="Student not found."
        className="m-8 w-auto"
      />
    );
  }

  const completedCount = todos.filter(
    (todo) => todo.completed_at != null,
  ).length;

  return (
    <div>
      <h2>Student Profile</h2>
      <ul>
        <li>Name: {user?.name ?? "Unknown"}</li>
        <li>Graduation Year: {graduationYearLabel(student.graduation_year)}</li>
        <li>High School: {student.organization ?? "Not set"}</li>
        <li>Major of Interest: {student.student_interest ?? "Not set"}</li>
        <li>
          Tasks Completed: {completedCount} / {todos.length}
        </li>
      </ul>
    </div>
  );
}
