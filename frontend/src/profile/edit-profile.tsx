import { useId, useState, type ReactNode } from "react";
import { Link, useNavigate } from "react-router-dom";
import {
  ArrowLeft,
  ChevronDown,
  CircleAlert,
  LogOut,
  PersonStanding,
} from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { cn } from "@/lib/utils";
import { errorMessage } from "@/api/client";
import {
  getGetCurrentUserQueryKey,
  getGetNotificationPreferencesQueryKey,
  getGetStudentQueryKey,
  useGetCurrentUser,
  useGetNotificationPreferences,
  useGetStudent,
  useUpdateNotificationPreferences,
  useUpdateStudent,
  useUpdateUser,
} from "@/api/endpoints";
import type {
  CurrentUser,
  NotificationPreferences,
  Student,
  UpdateNotificationPreferencesRequestBody,
  User,
  UpdateStudentBody,
} from "@/api/models";
import { ProfilePictureField } from "@/profile/components/pfp-dialog";
import { SaveButton, type SaveState } from "@/profile/components/save-button";
import UserAvatar from "@/components/user-avatar";
import { Skeleton } from "@/components/ui/skeleton";
import { DiscardChangesDialog } from "@/profile/components/discard-changes-dialog";
import { useLeaveGuard } from "@/profile/use-leave-guard";
import { useAuthContext } from "@/contexts/use-auth-context";
import { EmailValue } from "@/profile/components/email-field";
import { graduationYearOptions } from "@/lib/graduation-year";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";

const REGION_OPTIONS = [
  "Northeast",
  "Southeast",
  "Midwest",
  "Southwest",
  "West",
  "No Preference",
];

type FormValues = {
  name: string;
  highSchool: string;
  graduationYear: string;
  gpa: string;
  intendedMajor: string;
  intendedMinor: string;
  region: string;
};

const toFormValues = (user: CurrentUser, student?: Student): FormValues => ({
  name: user.name,
  highSchool: student?.high_school ?? "",
  graduationYear: student ? String(student.graduation_year) : "",
  gpa: student ? String(student.gpa) : "",
  intendedMajor: student?.intended_major ?? "",
  intendedMinor: student?.intended_minor ?? "",
  region: student?.geographic_preferences ?? "",
});

type NotificationValues = {
  email: boolean;
  weeklySummary: boolean;
  dueDate: boolean;
  daysBeforeDue: string;
  pastDue: boolean;
};

const toNotificationValues = (
  prefs: NotificationPreferences,
): NotificationValues => ({
  email: prefs.email_enabled,
  weeklySummary: prefs.weekly_summary_enabled,
  dueDate: prefs.due_date_notifications_enabled,
  daysBeforeDue: String(prefs.days_before_due),
  pastDue: prefs.notify_past_due,
});

const sameNotificationValues = (a: NotificationValues, b: NotificationValues) =>
  (Object.keys(a) as (keyof NotificationValues)[]).every((k) => a[k] === b[k]);

const daysLabel = (days: number) =>
  days === 0 ? "On the due date" : `${days} day${days === 1 ? "" : "s"} before`;

const onOff = (on: boolean) => (on ? "On" : "Off");

const fieldClass =
  "w-full rounded-lg border border-greyish-brown bg-white px-4 py-3 font-host text-sm leading-5 text-dark-greyish-brown outline-none transition-[border-color,box-shadow] focus-visible:border-dark-green focus-visible:ring-[3px] focus-visible:ring-dark-green/20 disabled:cursor-not-allowed disabled:opacity-60 aria-invalid:border-danger";

function Field({
  label,
  error,
  children,
}: {
  label: string;
  error?: string;
  children: (id: string) => ReactNode;
}) {
  const id = useId();
  return (
    <div className="flex flex-col gap-2">
      <label htmlFor={id} className="font-host text-sm leading-5 text-charcoal">
        {label}
      </label>
      {children(id)}
      {error && <p className="font-host text-xs text-danger">{error}</p>}
    </div>
  );
}

function SelectField({
  id,
  value,
  onChange,
  options,
  placeholder,
  disabled,
}: {
  id: string;
  value: string;
  onChange: (value: string) => void;
  options: { value: string; label: string }[];
  placeholder: string;
  disabled?: boolean;
}) {
  // Keep a saved value selectable even if it is no longer a listed option.
  const allOptions =
    value && !options.some((o) => o.value === value)
      ? [{ value, label: value }, ...options]
      : options;

  return (
    <div className="relative">
      <select
        id={id}
        data-slot="profile-select"
        value={value}
        disabled={disabled}
        onChange={(e) => onChange(e.target.value)}
        className={cn(fieldClass, "appearance-none pr-12")}
      >
        <option value="" disabled>
          {placeholder}
        </option>
        {allOptions.map((o) => (
          <option key={o.value} value={o.value}>
            {o.label}
          </option>
        ))}
      </select>
      <ChevronDown
        aria-hidden="true"
        className="pointer-events-none absolute top-1/2 right-4 size-6 -translate-y-1/2 text-charcoal"
        strokeWidth={1.5}
      />
    </div>
  );
}

function ErrorBubble({
  message,
  attempt,
  label = "Couldn't save changes",
}: {
  message: string;
  attempt: number;
  label?: string;
}) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <button
          type="button"
          data-slot="save-error"
          aria-label={label}
          className="flex size-8 cursor-help items-center justify-center rounded-full bg-danger/10 text-danger outline-none hover:bg-danger/15 focus-visible:ring-[3px] focus-visible:ring-danger/30"
        >
          {/* Remounting the icon restarts its animation. */}
          <CircleAlert
            key={attempt}
            className="save-error-shake size-5"
            strokeWidth={1.75}
          />
        </button>
      </TooltipTrigger>
      <TooltipContent side="bottom" sideOffset={6} className="max-w-xs">
        {message}
      </TooltipContent>
    </Tooltip>
  );
}

function Panel({
  title,
  aside,
  children,
}: {
  title: string;
  aside?: ReactNode;
  children: ReactNode;
}) {
  return (
    <section className="flex w-full flex-col gap-3 rounded-md border border-cream bg-white p-6 shadow-card">
      <h2 className="font-jakarta text-xl leading-[normal] font-medium text-charcoal">
        {title}
      </h2>
      <div className="flex flex-col gap-6 pt-2 sm:flex-row sm:items-center">
        {aside}
        <div className="grid min-w-0 flex-1 grid-cols-1 gap-4 sm:grid-cols-2">
          {children}
        </div>
      </div>
    </section>
  );
}

function SwitchRow({
  label,
  description,
  checked,
  onChange,
  disabled,
}: {
  label: string;
  description: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
}) {
  const id = useId();
  return (
    <div
      className={cn(
        "flex items-center justify-between gap-4 sm:col-span-2",
        disabled && "opacity-60",
      )}
    >
      <div className="flex min-w-0 flex-col gap-0.5">
        <label
          htmlFor={id}
          className="font-host text-sm leading-5 text-charcoal"
        >
          {label}
        </label>
        <p
          id={`${id}-desc`}
          className="font-host text-xs text-dark-greyish-brown"
        >
          {description}
        </p>
      </div>
      <button
        id={id}
        type="button"
        role="switch"
        data-slot="profile-switch"
        aria-checked={checked}
        aria-describedby={`${id}-desc`}
        disabled={disabled}
        onClick={() => onChange(!checked)}
        className={cn(
          "relative inline-flex h-6 w-11 shrink-0 items-center rounded-full outline-none transition-colors focus-visible:ring-[3px] focus-visible:ring-dark-green/30 disabled:cursor-not-allowed",
          checked ? "bg-dark-green" : "bg-greyish-brown",
        )}
      >
        <span
          aria-hidden="true"
          className={cn(
            "size-5 rounded-full bg-white shadow-sm transition-transform motion-reduce:transition-none",
            checked ? "translate-x-[22px]" : "translate-x-0.5",
          )}
        />
      </button>
    </div>
  );
}

function PageHeader({
  title,
  backLabel = "Go back",
  onBack,
  children,
}: {
  title: ReactNode;
  backLabel?: string;
  onBack?: () => void;
  children?: ReactNode;
}) {
  const navigate = useNavigate();
  return (
    <div className="flex w-full flex-wrap items-center justify-between gap-4">
      <div className="flex min-w-0 items-center gap-7">
        <button
          type="button"
          data-slot="back"
          aria-label={backLabel}
          onClick={onBack ?? (() => void navigate(-1))}
          className="flex size-8 shrink-0 items-center justify-center rounded-md text-charcoal outline-none hover:bg-black/5 focus-visible:ring-[3px] focus-visible:ring-dark-green/30"
        >
          <ArrowLeft className="size-8" strokeWidth={1.5} />
        </button>
        <h1 className="min-w-0 font-jakarta text-[30px] leading-[normal] font-bold break-words text-charcoal">
          {title}
        </h1>
      </div>
      <div className="flex items-center gap-2">{children}</div>
    </div>
  );
}

const EMPTY_VALUE = "—";

function ReadOnlyField({ label, value }: { label: string; value: ReactNode }) {
  const empty = value === null || value === undefined || value === "";
  return (
    <div className="flex min-w-0 flex-col gap-1">
      <dt className="font-host text-xs leading-4 text-dark-greyish-brown">
        {label}
      </dt>
      <dd
        className={cn(
          "font-host text-sm leading-5 break-words",
          empty ? "text-warm-grey" : "text-charcoal",
        )}
      >
        {empty ? EMPTY_VALUE : value}
      </dd>
    </div>
  );
}

function ReadOnlyPanel({
  title,
  aside,
  children,
}: {
  title: string;
  aside?: ReactNode;
  children: ReactNode;
}) {
  return (
    <section className="flex w-full flex-col gap-3 rounded-md border border-cream bg-white p-6 shadow-card">
      <h2 className="font-jakarta text-xl leading-[normal] font-medium text-charcoal">
        {title}
      </h2>
      <div className="flex flex-col gap-6 pt-2 sm:flex-row sm:items-center">
        {aside}
        <dl className="grid min-w-0 flex-1 grid-cols-1 gap-4 sm:grid-cols-2">
          {children}
        </dl>
      </div>
    </section>
  );
}

const editButtonClass =
  "flex h-8 items-center rounded-lg bg-dark-green px-4 font-host text-xs leading-[normal] whitespace-nowrap text-off-white outline-none transition-colors hover:bg-dark-green/90 focus-visible:ring-[3px] focus-visible:ring-dark-green/30";

const logOutButtonClass =
  "flex h-8 items-center gap-1.5 rounded-lg border border-dark-green bg-white px-4 font-host text-xs leading-[normal] whitespace-nowrap text-dark-green outline-none transition-colors hover:bg-dark-green/5 focus-visible:ring-[3px] focus-visible:ring-dark-green/30 disabled:cursor-not-allowed disabled:opacity-60";

function LogOutButton() {
  const { logout } = useAuthContext();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<unknown>(null);
  const [failures, setFailures] = useState(0);

  const handleLogOut = async () => {
    setPending(true);
    setError(null);
    try {
      await logout();
    } catch (err) {
      setError(err);
      setFailures((n) => n + 1);
      setPending(false);
    }
  };

  return (
    <>
      {error !== null && (
        <ErrorBubble
          message={errorMessage(error)}
          attempt={failures}
          label="Couldn't log out"
        />
      )}
      <p role="alert" className="sr-only">
        {error !== null && `Couldn't log out: ${errorMessage(error)}`}
      </p>
      <button
        type="button"
        data-slot="log-out"
        onClick={() => void handleLogOut()}
        disabled={pending}
        aria-busy={pending}
        className={logOutButtonClass}
      >
        <LogOut aria-hidden="true" className="size-3.5" strokeWidth={1.75} />
        {pending ? "Logging out…" : "Log Out"}
      </button>
    </>
  );
}

function ReadOnlyFieldSkeleton() {
  return (
    <div className="flex flex-col gap-1">
      <Skeleton className="h-4 w-24 bg-cream" />
      <Skeleton className="h-5 w-40 bg-cream" />
    </div>
  );
}

export function ProfileSkeleton({
  academic = true,
  own = true,
}: {
  academic?: boolean;
  own?: boolean;
}) {
  return (
    <PageShell>
      <PageHeader title={<Skeleton className="h-9 w-56 bg-cream" />}>
        {own && (
          <>
            <Skeleton className="h-8 w-[92px] rounded-lg bg-cream" />
            <Skeleton className="h-8 w-[98px] rounded-lg bg-cream" />
          </>
        )}
      </PageHeader>
      <div
        role="status"
        aria-label={own ? "Loading your profile" : "Loading profile"}
        className="flex w-full flex-col gap-6 sm:pl-15"
      >
        {[
          { title: "Personal Information", fields: 2, avatar: true },
          ...(academic
            ? [
                { title: "Academic Overview", fields: 3 },
                { title: "College Preferences & Specializations", fields: 3 },
              ]
            : []),
          ...(own ? [{ title: "Notifications", fields: 4 }] : []),
        ].map(({ title, fields, avatar }) => (
          <ReadOnlyPanel
            key={title}
            title={title}
            aside={
              avatar && (
                <Skeleton className="size-24 shrink-0 rounded-full bg-cream" />
              )
            }
          >
            {Array.from({ length: fields }, (_, i) => (
              <ReadOnlyFieldSkeleton key={i} />
            ))}
          </ReadOnlyPanel>
        ))}
      </div>
    </PageShell>
  );
}

export function PageShell({ children }: { children: ReactNode }) {
  return (
    <div className="flex min-h-full w-full flex-col bg-off-white py-12 pr-20 pl-12 max-lg:px-6">
      <div className="flex w-full max-w-[1090px] flex-col gap-6">
        {children}
      </div>
    </div>
  );
}

function ProfileForm({
  user,
  student,
  notifications,
  onCancel,
  onSaved,
}: {
  user: CurrentUser;
  student?: Student;
  notifications: NotificationPreferences;
  onCancel: () => void;
  onSaved: () => void;
}) {
  const [initial] = useState(() => toFormValues(user, student));
  const [values, setValues] = useState(initial);
  const [initialNotifications] = useState(() =>
    toNotificationValues(notifications),
  );
  const [notify, setNotify] = useState(initialNotifications);
  const [saveState, setSaveState] = useState<SaveState>("idle");
  const [failedSaves, setFailedSaves] = useState(0);

  const updateUser = useUpdateUser();
  const updateStudent = useUpdateStudent();
  const updateNotifications = useUpdateNotificationPreferences();

  const setNotification =
    <K extends keyof NotificationValues>(key: K) =>
    (value: NotificationValues[K]) =>
      setNotify((n) => ({ ...n, [key]: value }));

  const set = (key: keyof FormValues) => (value: string) => {
    setValues((v) => ({ ...v, [key]: value }));
  };

  const name = values.name.trim();
  const gpa = Number(values.gpa);
  const gpaError =
    student &&
    (values.gpa.trim() === "" || !Number.isFinite(gpa) || gpa < 0 || gpa > 5)
      ? "Enter a GPA between 0 and 5."
      : undefined;
  const nameError = name === "" ? "Enter your name." : undefined;
  const daysEditable = notify.email && notify.dueDate;
  const typedDays = Number(notify.daysBeforeDue);
  const daysError =
    daysEditable &&
    (notify.daysBeforeDue.trim() === "" ||
      !Number.isInteger(typedDays) ||
      typedDays < 0 ||
      typedDays > 30)
      ? "Enter a whole number of days from 0 to 30."
      : undefined;
  const daysBeforeDue = daysEditable
    ? typedDays
    : Number(initialNotifications.daysBeforeDue);

  const studentChanges = (): UpdateStudentBody => {
    if (!student) return {};
    const changes: UpdateStudentBody = {};
    if (values.highSchool !== initial.highSchool)
      changes.high_school = values.highSchool.trim();
    if (values.graduationYear !== initial.graduationYear)
      changes.graduation_year = Number(values.graduationYear);
    if (values.gpa !== initial.gpa) changes.gpa = gpa;
    if (values.intendedMajor !== initial.intendedMajor)
      changes.intended_major = values.intendedMajor.trim();
    if (values.intendedMinor !== initial.intendedMinor)
      changes.intended_minor = values.intendedMinor.trim();
    if (values.region !== initial.region)
      changes.geographic_preferences = values.region;
    return changes;
  };

  const nameChanged = values.name !== initial.name;
  const changes = studentChanges();
  const notificationsChanged = !sameNotificationValues(
    notify,
    initialNotifications,
  );
  const isDirty =
    nameChanged || Object.keys(changes).length > 0 || notificationsChanged;
  const hasErrors = !!nameError || !!gpaError || !!daysError;
  const leaveGuard = useLeaveGuard(isDirty);
  const error =
    updateUser.error ?? updateStudent.error ?? updateNotifications.error;

  const notificationBody = (): UpdateNotificationPreferencesRequestBody => ({
    email_enabled: notify.email,
    weekly_summary_enabled: notify.weeklySummary,
    due_date_notifications_enabled: notify.dueDate,
    days_before_due: daysBeforeDue,
    notify_past_due: notify.pastDue,
  });

  const save = async () => {
    if (hasErrors || !isDirty || saveState !== "idle") return;
    setSaveState("saving");
    try {
      await Promise.all([
        nameChanged && updateUser.mutateAsync({ data: { name } }),
        student &&
          Object.keys(changes).length > 0 &&
          updateStudent.mutateAsync({ id: student.id, data: changes }),
        notificationsChanged &&
          updateNotifications.mutateAsync({ data: notificationBody() }),
      ]);
      setSaveState("saved");
    } catch {
      setSaveState("idle");
      setFailedSaves((n) => n + 1);
    }
  };

  return (
    <PageShell>
      <form
        className="contents"
        onSubmit={(e) => {
          e.preventDefault();
          void save();
        }}
      >
        <PageHeader
          title="Edit Profile"
          backLabel="Back to profile"
          onBack={() => leaveGuard.guard(onCancel)}
        >
          {error && (
            <ErrorBubble message={errorMessage(error)} attempt={failedSaves} />
          )}
          <SaveButton
            state={saveState}
            disabled={!isDirty || hasErrors}
            onSavedShown={onSaved}
          />
        </PageHeader>

        <p role="alert" className="sr-only">
          {error && `Couldn't save: ${errorMessage(error)}`}
        </p>

        <div className="flex w-full flex-col gap-6 sm:pl-15">
          <Panel
            title="Personal Information"
            aside={<ProfilePictureField user={user} />}
          >
            <Field label="Name" error={nameError}>
              {(id) => (
                <input
                  id={id}
                  data-slot="profile-input"
                  className={fieldClass}
                  value={values.name}
                  onChange={(e) => set("name")(e.target.value)}
                  autoComplete="name"
                  aria-invalid={!!nameError}
                />
              )}
            </Field>
            {student && (
              <>
                <Field label="High School">
                  {(id) => (
                    <input
                      id={id}
                      data-slot="profile-input"
                      className={fieldClass}
                      value={values.highSchool}
                      onChange={(e) => set("highSchool")(e.target.value)}
                    />
                  )}
                </Field>
                <Field label="Graduation Year">
                  {(id) => (
                    <SelectField
                      id={id}
                      value={values.graduationYear}
                      onChange={set("graduationYear")}
                      placeholder="Select year"
                      options={graduationYearOptions().map(
                        ({ value, label }) => ({ value: String(value), label }),
                      )}
                    />
                  )}
                </Field>
                <Field label="GPA (Unweighted)" error={gpaError}>
                  {(id) => (
                    <input
                      id={id}
                      data-slot="profile-input"
                      className={fieldClass}
                      inputMode="decimal"
                      value={values.gpa}
                      onChange={(e) => set("gpa")(e.target.value)}
                      aria-invalid={!!gpaError}
                    />
                  )}
                </Field>
              </>
            )}
          </Panel>

          {student && (
            <Panel title="College Preferences & Specializations">
              <Field label="Intended Major">
                {(id) => (
                  <input
                    id={id}
                    data-slot="profile-input"
                    className={fieldClass}
                    value={values.intendedMajor}
                    onChange={(e) => set("intendedMajor")(e.target.value)}
                  />
                )}
              </Field>
              <Field label="Minor / Secondary Academic Interests">
                {(id) => (
                  <input
                    id={id}
                    data-slot="profile-input"
                    className={fieldClass}
                    value={values.intendedMinor}
                    onChange={(e) => set("intendedMinor")(e.target.value)}
                  />
                )}
              </Field>
              <Field label="Geographic Region Preference">
                {(id) => (
                  <SelectField
                    id={id}
                    value={values.region}
                    onChange={set("region")}
                    placeholder="Select region"
                    options={REGION_OPTIONS.map((r) => ({
                      value: r,
                      label: r,
                    }))}
                  />
                )}
              </Field>
            </Panel>
          )}

          <Panel title="Notifications">
            <SwitchRow
              label="Email notifications"
              description="Turn off to stop all notification emails."
              checked={notify.email}
              onChange={setNotification("email")}
            />
            <SwitchRow
              label="Weekly summary"
              description="A weekly email of your outstanding assignments."
              checked={notify.weeklySummary}
              onChange={setNotification("weeklySummary")}
              disabled={!notify.email}
            />
            <SwitchRow
              label="Due date reminders"
              description="An email before an assignment is due."
              checked={notify.dueDate}
              onChange={setNotification("dueDate")}
              disabled={!notify.email}
            />
            <Field label="Days before due date" error={daysError}>
              {(id) => (
                <input
                  id={id}
                  data-slot="profile-input"
                  className={fieldClass}
                  type="number"
                  inputMode="numeric"
                  min={0}
                  max={30}
                  step={1}
                  value={notify.daysBeforeDue}
                  onChange={(e) =>
                    setNotification("daysBeforeDue")(e.target.value)
                  }
                  disabled={!daysEditable}
                  aria-invalid={!!daysError}
                />
              )}
            </Field>
            <SwitchRow
              label="Past due alerts"
              description="An email when an assignment becomes past due."
              checked={notify.pastDue}
              onChange={setNotification("pastDue")}
              disabled={!notify.email}
            />
          </Panel>
        </div>
      </form>

      <DiscardChangesDialog {...leaveGuard.dialog} />
    </PageShell>
  );
}

type ProfileUser = Pick<User, "id" | "name" | "email" | "pfp_key"> & {
  pfp_url?: string | null;
};

interface OwnProfile {
  notifications: NotificationPreferences;
  onEdit: () => void;
  canInvite?: boolean;
}

export function ProfileView({
  user,
  student,
  own,
}: {
  user: ProfileUser;
  student?: Student;
  own?: OwnProfile;
}) {
  return (
    <PageShell>
      <PageHeader title={user.name}>
        {own && (
          <>
            <LogOutButton />
            <button
              type="button"
              data-slot="edit-profile"
              onClick={own.onEdit}
              className={editButtonClass}
            >
              Edit Profile
            </button>
          </>
        )}
      </PageHeader>

      <div className="flex w-full flex-col gap-6 sm:pl-15">
        <ReadOnlyPanel
          title="Personal Information"
          aside={
            <UserAvatar
              name={user.name}
              userId={user.id}
              pfpKey={user.pfp_key}
              src={user.pfp_url}
              className="size-24 shrink-0"
              fallbackClassName="bg-mist text-3xl text-neutral-600"
            />
          }
        >
          <ReadOnlyField label="Name" value={user.name} />
          <ReadOnlyField
            label="Email"
            value={
              own || !user.email ? (
                user.email
              ) : (
                <EmailValue email={user.email} userId={user.id} />
              )
            }
          />
        </ReadOnlyPanel>

        {student && (
          <>
            <ReadOnlyPanel title="Academic Overview">
              <ReadOnlyField label="High School" value={student.high_school} />
              <ReadOnlyField
                label="Graduation Year"
                value={student.graduation_year}
              />
              <ReadOnlyField label="GPA (Unweighted)" value={student.gpa} />
            </ReadOnlyPanel>

            <ReadOnlyPanel title="College Preferences & Specializations">
              <ReadOnlyField
                label="Intended Major"
                value={student.intended_major}
              />
              <ReadOnlyField
                label="Minor / Secondary Academic Interests"
                value={student.intended_minor}
              />
              <ReadOnlyField
                label="Geographic Region Preference"
                value={student.geographic_preferences}
              />
            </ReadOnlyPanel>
          </>
        )}

        {own && <NotificationsPanel notifications={own.notifications} />}

        {own?.canInvite && <InviteUserCard />}
      </div>
    </PageShell>
  );
}

const INVITE_USER_PATH = "/create-account";

function InviteUserCard() {
  return (
    <section
      aria-label="Invite a user"
      className="flex w-full flex-wrap items-center justify-between gap-4 rounded-md border border-cream bg-white p-6 shadow-card"
    >
      <div className="flex items-center gap-4">
        <PersonStanding
          aria-hidden="true"
          className="size-11 shrink-0 text-charcoal"
          strokeWidth={1.5}
        />
        <div className="flex flex-col">
          <p className="font-host text-sm leading-5 text-charcoal">
            Don&apos;t see a student in the directory?
          </p>
          <p className="font-host text-xs leading-4 text-dark-greyish-brown">
            Invite a user to the Inspirate Consulting Hub!
          </p>
        </div>
      </div>
      <Link
        to={INVITE_USER_PATH}
        data-slot="send-invite"
        className="flex h-9 items-center rounded-lg border border-greyish-brown bg-white px-4 font-host text-xs leading-[normal] whitespace-nowrap text-charcoal outline-none transition-colors hover:bg-black/5 focus-visible:ring-[3px] focus-visible:ring-dark-green/30"
      >
        Send Invite
      </Link>
    </section>
  );
}

function NotificationsPanel({
  notifications,
}: {
  notifications: NotificationPreferences;
}) {
  const emailOn = notifications.email_enabled;
  const sub = (on: boolean) => (emailOn ? onOff(on) : "Off (email is off)");
  return (
    <ReadOnlyPanel title="Notifications">
      <ReadOnlyField label="Email notifications" value={onOff(emailOn)} />
      <ReadOnlyField
        label="Weekly summary"
        value={sub(notifications.weekly_summary_enabled)}
      />
      <ReadOnlyField
        label="Due date reminders"
        value={
          emailOn && notifications.due_date_notifications_enabled
            ? `On, ${daysLabel(notifications.days_before_due).toLowerCase()}`
            : sub(notifications.due_date_notifications_enabled)
        }
      />
      <ReadOnlyField
        label="Past due alerts"
        value={sub(notifications.notify_past_due)}
      />
    </ReadOnlyPanel>
  );
}

function Profile({
  user,
  student,
  notifications,
}: {
  user: CurrentUser;
  student?: Student;
  notifications: NotificationPreferences;
}) {
  const queryClient = useQueryClient();
  const [editing, setEditing] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [announcement, setAnnouncement] = useState("");

  const handleSaved = () => {
    setEditing(false);
    setRefreshing(true);
    setAnnouncement("Changes saved.");
    void Promise.all([
      queryClient.invalidateQueries({ queryKey: getGetCurrentUserQueryKey() }),
      queryClient.invalidateQueries({
        queryKey: getGetNotificationPreferencesQueryKey(),
      }),
      student &&
        queryClient.invalidateQueries({
          queryKey: getGetStudentQueryKey(student.id),
        }),
    ]).finally(() => setRefreshing(false));
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
        <ProfileView
          user={user}
          student={student}
          own={{
            notifications,
            canInvite: user.role === "counselor",
            onEdit: () => {
              setAnnouncement("");
              setEditing(true);
            },
          }}
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
      notifications={notifications.data}
    />
  );
}
