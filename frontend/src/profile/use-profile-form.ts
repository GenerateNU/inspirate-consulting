import { useState } from "react";
import type {
  Counselor,
  CurrentUser,
  NotificationPreferences,
  Student,
} from "@/api/models";
import {
  counselorChanges,
  sameNotificationValues,
  studentChanges,
  toNotificationBody,
  toNotificationValues,
  toProfileValues,
  validateProfile,
  type SetNotificationValue,
  type SetProfileValue,
} from "@/profile/profile-form-values";
import { useSaveProfile, type ProfileUpdate } from "@/profile/use-save-profile";

export interface ProfileSource {
  user: CurrentUser;
  student?: Student;
  counselor?: Counselor;
  notifications: NotificationPreferences;
}

export function useProfileForm({
  user,
  student,
  counselor,
  notifications,
}: ProfileSource) {
  const [initial] = useState(() => toProfileValues(user, student, counselor));
  const [values, setValues] = useState(initial);
  const [initialNotify] = useState(() => toNotificationValues(notifications));
  const [notify, setNotify] = useState(initialNotify);
  const saver = useSaveProfile();

  const setValue: SetProfileValue = (key) => (value) =>
    setValues((v) => ({ ...v, [key]: value }));
  const setNotification: SetNotificationValue = (key) => (value) =>
    setNotify((n) => ({ ...n, [key]: value }));

  const changes = student ? studentChanges(initial, values) : {};
  const counselorUpdate = counselor ? counselorChanges(initial, values) : {};
  const update: ProfileUpdate = {
    name: values.name !== initial.name ? values.name.trim() : undefined,
    student:
      student && Object.keys(changes).length > 0
        ? { id: student.id, changes }
        : undefined,
    counselor:
      Object.keys(counselorUpdate).length > 0 ? counselorUpdate : undefined,
    notifications: sameNotificationValues(notify, initialNotify)
      ? undefined
      : toNotificationBody(notify, initialNotify),
  };
  const isDirty = Object.values(update).some((part) => part !== undefined);
  const errors = validateProfile(
    initial,
    values,
    notify,
    !!student,
    !!counselor,
  );
  const hasErrors = Object.values(errors).some(Boolean);

  const save = async () => {
    if (hasErrors || !isDirty) return;
    await saver.save(update);
  };

  return {
    values,
    setValue,
    notify,
    setNotification,
    errors,
    hasErrors,
    isDirty,
    save,
    saveState: saver.saveState,
    failedSaves: saver.failedSaves,
    error: saver.error,
  };
}
