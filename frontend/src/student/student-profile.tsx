// name, grade, highschool, major of interest, tasks completed
import { useParams } from "react-router-dom";
import { useGetStudent, useFetchUser, useGetTodoItems } from "../api/endpoints";

const gradeByYear: Record<string, number> = {
  freshman: 9,
  sophomore: 10,
  junior: 11,
  senior: 12,
};

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

  if (studentQuery.error) {
    return <p>Could not load student: {studentQuery.error.message}</p>;
  }

  if (userQuery.error) {
    return <p>Could not load user: {userQuery.error.message}</p>;
  }

  if (todosQuery.error) {
    return <p>Could not load tasks: {todosQuery.error.message}</p>;
  }

  if (!student) {
    return <p>Student not found.</p>;
  }

  const completedCount = todos.filter(
    (todo) => todo.completed_at != null,
  ).length;

  return (
    <div>
      <h2>Student Profile</h2>
      <ul>
        <li>Name: {user?.name ?? "Unknown"}</li>
        <li>Grade: {gradeByYear[student.year] ?? student.year}</li>
        <li>High School: {student.organization ?? "Not set"}</li>
        <li>Major of Interest: {student.student_interest ?? "Not set"}</li>
        <li>
          Tasks Completed: {completedCount} / {todos.length}
        </li>
      </ul>
    </div>
  );
}
