import { useState } from "react";
import {
  useUpdateNotificationPreferences,
  useUpdateStudent,
  useUpdateUser,
} from "@/api/endpoints";
import type {
  UpdateNotificationPreferencesRequestBody,
  UpdateStudentBody,
} from "@/api/models";
import type { SaveState } from "@/profile/components/save-button";

export interface ProfileUpdate {
  name?: string;
  student?: StudentUpdate;
  notifications?: UpdateNotificationPreferencesRequestBody;
}

export interface StudentUpdate {
  id: string;
  changes: UpdateStudentBody;
}

export function useSaveProfile() {
  const [saveState, setSaveState] = useState<SaveState>("idle");
  const [failedSaves, setFailedSaves] = useState(0);
  const updateUser = useUpdateUser();
  const updateStudent = useUpdateStudent();
  const updateNotifications = useUpdateNotificationPreferences();

  const send = ({ name, student, notifications }: ProfileUpdate) =>
    Promise.all([
      name !== undefined && updateUser.mutateAsync({ data: { name } }),
      student &&
        updateStudent.mutateAsync({ id: student.id, data: student.changes }),
      notifications && updateNotifications.mutateAsync({ data: notifications }),
    ]);

  const save = async (update: ProfileUpdate) => {
    if (saveState !== "idle") return;
    setSaveState("saving");
    try {
      await send(update);
      setSaveState("saved");
    } catch {
      setSaveState("idle");
      setFailedSaves((n) => n + 1);
    }
  };

  const error =
    updateUser.error ?? updateStudent.error ?? updateNotifications.error;

  return { save, saveState, failedSaves, error };
}
