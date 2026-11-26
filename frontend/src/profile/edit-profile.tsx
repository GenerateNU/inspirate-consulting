import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import {
  getGetCurrentUserQueryKey,
  getGetNotificationPreferencesQueryKey,
  getGetStudentQueryKey,
  useGetCurrentUser,
  useGetNotificationPreferences,
  useGetStudent,
} from "@/api/endpoints";
import { PageShell } from "@/profile/components/page-shell";
import { ProfileForm } from "@/profile/components/profile-form";
import { ProfileSkeleton } from "@/profile/components/profile-skeleton";
import { ProfileView } from "@/profile/components/profile-view";
import type { ProfileSource } from "@/profile/use-profile-form";
import { ErrorMessage } from "@/components/ui/error-message";

function useRefetchProfile(studentId?: string) {
  const queryClient = useQueryClient();
  return () =>
    Promise.all([
      queryClient.invalidateQueries({ queryKey: getGetCurrentUserQueryKey() }),
      queryClient.invalidateQueries({
        queryKey: getGetNotificationPreferencesQueryKey(),
      }),
      studentId &&
        queryClient.invalidateQueries({
          queryKey: getGetStudentQueryKey(studentId),
        }),
    ]);
}

function Profile({ user, student, notifications }: ProfileSource) {
  const refetch = useRefetchProfile(student?.id);
  const [editing, setEditing] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [announcement, setAnnouncement] = useState("");

  const handleSaved = () => {
    setEditing(false);
    setRefreshing(true);
    setAnnouncement("Changes saved.");
    void refetch().finally(() => setRefreshing(false));
  };

  const handleEdit = () => {
    setAnnouncement("");
    setEditing(true);
  };

  const own = {
    notifications,
    canInvite: user.role === "counselor",
    onEdit: handleEdit,
  };

  return (
    <>
      <p role="status" className="sr-only">
        {announcement}
      </p>
      {editing ? (
        <ProfileForm
          user={user}
          student={student}
          notifications={notifications}
          onCancel={() => setEditing(false)}
          onSaved={handleSaved}
        />
      ) : refreshing ? (
        <ProfileSkeleton academic={!!student} />
      ) : (
        <ProfileView user={user} student={student} own={own} />
      )}
    </>
  );
}

export default function EditProfile() {
  const currentUser = useGetCurrentUser();
  const studentID = currentUser.data?.student_id ?? "";
  const student = useGetStudent(studentID, {
    query: { enabled: studentID !== "" },
  });
  const notifications = useGetNotificationPreferences();

  const waitingForStudent = studentID !== "" && !student.data && !student.error;
  if (currentUser.isLoading || waitingForStudent || notifications.isLoading) {
    return (
      <ProfileSkeleton academic={currentUser.data?.role !== "counselor"} />
    );
  }

  const loadError = currentUser.error ?? student.error ?? notifications.error;
  if (loadError || !currentUser.data || !notifications.data) {
    return (
      <PageShell>
        <ErrorMessage
          variant="banner"
          title="Couldn't load your profile"
          error={loadError}
          message={loadError ? undefined : "Something went wrong."}
          onRetry={() =>
            [currentUser, student, notifications]
              .filter((query) => query.error)
              .forEach((query) => void query.refetch())
          }
        />
      </PageShell>
    );
  }

  return (
    <Profile
      user={currentUser.data}
      student={student.data}
      notifications={notifications.data}
    />
  );
}
