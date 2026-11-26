import { useFetchUser, useGetStudent } from "../../api/endpoints";

// one student plus their name, which lives on the user row
export function useStudentProfile(studentId: string) {
  const studentQuery = useGetStudent(studentId, {
    query: { enabled: studentId !== "" },
  });
  const student = studentQuery.data;

  const userId = student?.user_id ?? "";
  const userQuery = useFetchUser(userId, { query: { enabled: userId !== "" } });

  // "" until the name has loaded, so headers don't flash "Student"; "Student"
  // only when it can't be loaded at all
  let name = "";
  if (userQuery.data) name = userQuery.data.name;
  else if (studentQuery.error || userQuery.error) name = "Student";

  return {
    student,
    name,
    isLoading: studentQuery.isLoading || userQuery.isLoading,
    error: studentQuery.error,
    refetch: studentQuery.refetch,
  };
}
