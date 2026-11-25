import { Navigate, useParams } from "react-router-dom";
import { errorMessage } from "@/api/client";
import { useFetchStudents, useFetchUser, useGetStudent } from "@/api/endpoints";
import { useAuthContext } from "@/contexts/use-auth-context";
import { PageShell } from "@/profile/components/page-shell";
import { ProfileSkeleton } from "@/profile/components/profile-skeleton";
import { ProfileView } from "@/profile/components/profile-view";

export default function ViewProfile() {
  const { userId = "" } = useParams();
  const { userID: currentUserId, role } = useAuthContext();
  const isSelf = userId === currentUserId;

  const user = useFetchUser(userId, {
    query: { enabled: userId !== "" && !isSelf },
  });

  const studentLookup = useFetchStudents(
    { user_id: userId, limit: 1 },
    { query: { enabled: userId !== "" && !isSelf && role === "counselor" } },
  );
  const studentId = studentLookup.data?.students?.[0]?.id ?? "";
  const student = useGetStudent(studentId, {
    query: { enabled: studentId !== "" },
  });

  if (isSelf) return <Navigate to="/profile" replace />;

  const waitingForStudent = studentId !== "" && !student.data && !student.error;
  if (user.isLoading || studentLookup.isLoading || waitingForStudent) {
    return <ProfileSkeleton academic={studentId !== ""} own={false} />;
  }

  const loadError = user.error ?? studentLookup.error ?? student.error;
  if (loadError || !user.data) {
    return (
      <PageShell>
        <p className="font-host text-sm text-danger">
          Could not load this profile: {errorMessage(loadError)}
        </p>
      </PageShell>
    );
  }

  return <ProfileView user={user.data} student={student.data} />;
}
