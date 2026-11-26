import type {
  Counselor,
  CurrentUser,
  NotificationPreferences,
  Student,
  UpdateCounselorBody,
  UpdateNotificationPreferencesRequestBody,
  UpdateStudentBody,
} from "@/api/models";

export interface ProfileValues {
  name: string;
  highSchool: string;
  graduationYear: string;
  gpa: string;
  intendedMajor: string;
  intendedMinor: string;
  region: string;
  customLink: string;
  expiresAt: string;
}

export interface NotificationValues {
  email: boolean;
  weeklySummary: boolean;
  dueDate: boolean;
  daysBeforeDue: string;
  pastDue: boolean;
}

export interface ProfileErrors {
  name?: string;
  gpa?: string;
  daysBeforeDue?: string;
  customLink?: string;
}

export type SetProfileValue = (
  key: keyof ProfileValues,
) => (value: string) => void;

export type SetNotificationValue = <K extends keyof NotificationValues>(
  key: K,
) => (value: NotificationValues[K]) => void;

// <input type="date"> works in local YYYY-MM-DD, the API in ISO timestamps
export function toDateInput(date: string | Date | null | undefined): string {
  if (!date) return "";
  const d = new Date(date);
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;
}

export function tomorrowDateInput(): string {
  const d = new Date();
  d.setDate(d.getDate() + 1);
  return toDateInput(d);
}

export const toProfileValues = (
  user: CurrentUser,
  student?: Student,
  counselor?: Counselor,
): ProfileValues => ({
  name: user.name,
  highSchool: student?.high_school ?? "",
  graduationYear: student ? String(student.graduation_year) : "",
  gpa: student ? String(student.gpa) : "",
  intendedMajor: student?.intended_major ?? "",
  intendedMinor: student?.intended_minor ?? "",
  region: student?.geographic_preferences ?? "",
  customLink: counselor?.custom_link ?? "",
  expiresAt: toDateInput(counselor?.expires_at),
});

export const toNotificationValues = (
  prefs: NotificationPreferences,
): NotificationValues => ({
  email: prefs.email_enabled,
  weeklySummary: prefs.weekly_summary_enabled,
  dueDate: prefs.due_date_notifications_enabled,
  daysBeforeDue: String(prefs.days_before_due),
  pastDue: prefs.notify_past_due,
});

export const sameNotificationValues = (
  a: NotificationValues,
  b: NotificationValues,
) =>
  (Object.keys(a) as (keyof NotificationValues)[]).every((k) => a[k] === b[k]);

export const daysEditable = (notify: NotificationValues) =>
  notify.email && notify.dueDate;

const isNumberInRange = (text: string, min: number, max: number) => {
  const n = Number(text);
  return text.trim() !== "" && Number.isFinite(n) && n >= min && n <= max;
};

const isWholeNumberInRange = (text: string, min: number, max: number) =>
  isNumberInRange(text, min, max) && Number.isInteger(Number(text));

const isWebUrl = (text: string) => {
  try {
    const { protocol } = new URL(text);
    return protocol === "http:" || protocol === "https:";
  } catch {
    return false;
  }
};

export function validateProfile(
  values: ProfileValues,
  notify: NotificationValues,
  hasStudent: boolean,
  hasCounselor: boolean,
): ProfileErrors {
  const errors: ProfileErrors = {};
  if (values.name.trim() === "") errors.name = "Enter your name.";
  if (hasStudent && !isNumberInRange(values.gpa, 0, 5))
    errors.gpa = "Enter a GPA between 0 and 5.";
  if (
    daysEditable(notify) &&
    !isWholeNumberInRange(notify.daysBeforeDue, 0, 30)
  )
    errors.daysBeforeDue = "Enter a whole number of days from 0 to 30.";
  if (hasCounselor) {
    const link = values.customLink.trim();
    if (link !== "" && !isWebUrl(link))
      errors.customLink = "Enter a full link starting with https://.";
  }
  return errors;
}

// a new expiration of today or earlier ends the account as soon as it's saved
export const expiresImmediately = (
  initial: ProfileValues,
  values: ProfileValues,
) =>
  values.expiresAt !== initial.expiresAt &&
  values.expiresAt !== "" &&
  values.expiresAt < tomorrowDateInput();

const toStudentBody = (values: ProfileValues): UpdateStudentBody => ({
  high_school: values.highSchool.trim(),
  graduation_year: Number(values.graduationYear),
  gpa: Number(values.gpa),
  intended_major: values.intendedMajor.trim(),
  intended_minor: values.intendedMinor.trim(),
  geographic_preferences: values.region,
});

export function studentChanges(
  initial: ProfileValues,
  values: ProfileValues,
): UpdateStudentBody {
  const before = toStudentBody(initial);
  const after = toStudentBody(values);
  return Object.fromEntries(
    Object.entries(after).filter(
      ([key, value]) => before[key as keyof UpdateStudentBody] !== value,
    ),
  );
}

// an emptied field is sent as null so the backend clears it
export function counselorChanges(
  initial: ProfileValues,
  values: ProfileValues,
): UpdateCounselorBody {
  const changes: UpdateCounselorBody = {};
  const link = values.customLink.trim();
  if (link !== initial.customLink) changes.custom_link = link || null;
  if (values.expiresAt !== initial.expiresAt)
    changes.expires_at = values.expiresAt
      ? new Date(`${values.expiresAt}T00:00:00`).toISOString()
      : null;
  return changes;
}

export const toNotificationBody = (
  notify: NotificationValues,
  saved: NotificationValues,
): UpdateNotificationPreferencesRequestBody => ({
  email_enabled: notify.email,
  weekly_summary_enabled: notify.weeklySummary,
  due_date_notifications_enabled: notify.dueDate,
  days_before_due: Number(
    daysEditable(notify) ? notify.daysBeforeDue : saved.daysBeforeDue,
  ),
  notify_past_due: notify.pastDue,
});
