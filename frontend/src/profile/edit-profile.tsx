import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { errorMessage } from "@/api/client";
import {
  getGetCounselorQueryKey,
  getGetCurrentUserQueryKey,
  getGetNotificationPreferencesQueryKey,
  getGetStudentQueryKey,
  useGetCounselor,
  useGetCurrentUser,
  useGetNotificationPreferences,
  useGetStudent,
} from "@/api/endpoints";
import { PageShell } from "@/profile/components/page-shell";
import { ProfileForm } from "@/profile/components/profile-form";
import { ProfileSkeleton } from "@/profile/components/profile-skeleton";
import { ProfileView } from "@/profile/components/profile-view";
import type { ProfileSource } from "@/profile/use-profile-form";

function useRefetchProfile(studentId?: string, counselorId?: string) {
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
      counselorId &&
        queryClient.invalidateQueries({
          queryKey: getGetCounselorQueryKey(counselorId),
        }),
    ]);
}

function Profile({ user, student, counselor, notifications }: ProfileSource) {
  const refetch = useRefetchProfile(student?.id, counselor?.id);
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
          counselor={counselor}
          notifications={notifications}
          onCancel={() => setEditing(false)}
          onSaved={handleSaved}
        />
      ) : refreshing ? (
        <ProfileSkeleton academic={!!student} counselor={!!counselor} />
      ) : (
        <ProfileView
          user={user}
          student={student}
          counselor={counselor}
          own={own}
        />
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
  const counselorID = currentUser.data?.counselor_id ?? "";
  const counselor = useGetCounselor(counselorID, {
    query: { enabled: counselorID !== "" },
  });
  const notifications = useGetNotificationPreferences();

  const waitingForStudent = studentID !== "" && !student.data && !student.error;
  const waitingForCounselor =
    counselorID !== "" && !counselor.data && !counselor.error;
  if (
    currentUser.isLoading ||
    waitingForStudent ||
    waitingForCounselor ||
    notifications.isLoading
  ) {
    const isCounselor = currentUser.data?.role === "counselor";
    return <ProfileSkeleton academic={!isCounselor} counselor={isCounselor} />;
  }

  const loadError =
    currentUser.error ??
    student.error ??
    counselor.error ??
    notifications.error;
  if (loadError || !currentUser.data || !notifications.data) {
    return (
      <PageShell>
        <p className="font-host text-sm text-danger">
          Could not load your profile: {errorMessage(loadError)}
        </p>
      </PageShell>
    );
  }

  return (
    <Profile
      user={currentUser.data}
      student={student.data}
      counselor={counselor.data}
      notifications={notifications.data}
    />
  );
}
