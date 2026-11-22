import { useId, useState, type ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import { ArrowLeft, ChevronDown } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { cn } from "@/lib/utils";
import { errorMessage } from "@/api/client";
import {
  getGetCurrentUserQueryKey,
  getGetStudentQueryKey,
  useGetCurrentUser,
  useGetStudent,
  useUpdateStudent,
  useUpdateUser,
} from "@/api/endpoints";
import {
  UpdateStudentBodyYear,
  type CurrentUser,
  type Student,
  type UpdateStudentBody,
} from "@/api/models";

const REGION_OPTIONS = [
  "Northeast",
  "Southeast",
  "Midwest",
  "Southwest",
  "West",
  "No Preference",
];

const CAMPUS_OPTIONS = ["City", "Suburban", "Rural", "No Preference"];

// Years of high school left after the current one, by grade.
const YEARS_LEFT: Record<UpdateStudentBodyYear, number> = {
  senior: 0,
  junior: 1,
  sophomore: 2,
  freshman: 3,
};

// The school year ends in the spring; from July on, seniors graduate next year.
const currentClassYear = () => {
  const now = new Date();
  return now.getMonth() >= 6 ? now.getFullYear() + 1 : now.getFullYear();
};

const graduationOptions = () =>
  (Object.keys(YEARS_LEFT) as UpdateStudentBodyYear[])
    .map((year) => ({
      year,
      label: String(currentClassYear() + YEARS_LEFT[year]),
    }))
    .sort((a, b) => a.label.localeCompare(b.label));

type FormValues = {
  firstName: string;
  lastName: string;
  highSchool: string;
  year: string;
  gpa: string;
  intendedMajor: string;
  intendedMinor: string;
  region: string;
};

const splitName = (name: string) => {
  const trimmed = name.trim();
  const space = trimmed.indexOf(" ");
  return space === -1
    ? { firstName: trimmed, lastName: "" }
    : {
        firstName: trimmed.slice(0, space),
        lastName: trimmed.slice(space + 1).trim(),
      };
};

const toFormValues = (user: CurrentUser, student?: Student): FormValues => ({
  ...splitName(user.name),
  highSchool: student?.high_school ?? "",
  year: student?.year ?? "",
  gpa: student ? String(student.gpa) : "",
  intendedMajor: student?.intended_major ?? "",
  intendedMinor: student?.intended_minor ?? "",
  region: student?.geographic_preferences ?? "",
});

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

function Panel({ title, children }: { title: string; children: ReactNode }) {
  return (
    <section className="flex w-full flex-col gap-3 rounded-md border border-cream bg-white p-6 shadow-card">
      <h2 className="font-jakarta text-xl leading-[normal] font-medium text-charcoal">
        {title}
      </h2>
      <div className="grid grid-cols-1 gap-4 pt-2 sm:grid-cols-2">
        {children}
      </div>
    </section>
  );
}

function PageShell({ children }: { children: ReactNode }) {
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
}: {
  user: CurrentUser;
  student?: Student;
}) {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [initial, setInitial] = useState(() => toFormValues(user, student));
  const [values, setValues] = useState(initial);
  const [status, setStatus] = useState("");

  const updateUser = useUpdateUser();
  const updateStudent = useUpdateStudent();

  const set = (key: keyof FormValues) => (value: string) => {
    setStatus("");
    setValues((v) => ({ ...v, [key]: value }));
  };

  const fullName =
    `${values.firstName.trim()} ${values.lastName.trim()}`.trim();
  const gpa = Number(values.gpa);
  const gpaError =
    student && (values.gpa.trim() === "" || !Number.isInteger(gpa) || gpa < 0)
      ? "Enter a whole number."
      : undefined;
  const nameError = fullName === "" ? "Enter your first name." : undefined;

  const studentChanges = (): UpdateStudentBody => {
    if (!student) return {};
    const changes: UpdateStudentBody = {};
    if (values.highSchool !== initial.highSchool)
      changes.high_school = values.highSchool.trim();
    if (values.year !== initial.year)
      changes.year = values.year as UpdateStudentBodyYear;
    if (values.gpa !== initial.gpa) changes.gpa = gpa;
    if (values.intendedMajor !== initial.intendedMajor)
      changes.intended_major = values.intendedMajor.trim();
    if (values.intendedMinor !== initial.intendedMinor)
      changes.intended_minor = values.intendedMinor.trim();
    if (values.region !== initial.region)
      changes.geographic_preferences = values.region;
    return changes;
  };

  const nameChanged =
    values.firstName !== initial.firstName ||
    values.lastName !== initial.lastName;
  const changes = studentChanges();
  const isDirty = nameChanged || Object.keys(changes).length > 0;
  const isSaving = updateUser.isPending || updateStudent.isPending;
  const error = updateUser.error ?? updateStudent.error;

  const save = async () => {
    if (nameError || gpaError || !isDirty) return;
    setStatus("");
    try {
      await Promise.all([
        nameChanged && updateUser.mutateAsync({ data: { name: fullName } }),
        student &&
          Object.keys(changes).length > 0 &&
          updateStudent.mutateAsync({ id: student.id, data: changes }),
      ]);
      await Promise.all([
        queryClient.invalidateQueries({
          queryKey: getGetCurrentUserQueryKey(),
        }),
        student &&
          queryClient.invalidateQueries({
            queryKey: getGetStudentQueryKey(student.id),
          }),
      ]);
      setInitial(values);
      setStatus("Changes saved.");
    } catch {
      // Shown from the mutation error below.
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
        <div className="flex w-full flex-wrap items-center justify-between gap-4">
          <div className="flex items-center gap-7">
            <button
              type="button"
              data-slot="back"
              aria-label="Go back"
              onClick={() => void navigate(-1)}
              className="flex size-8 items-center justify-center rounded-md text-charcoal outline-none hover:bg-black/5 focus-visible:ring-[3px] focus-visible:ring-dark-green/30"
            >
              <ArrowLeft className="size-8" strokeWidth={1.5} />
            </button>
            <h1 className="font-jakarta text-[30px] leading-[normal] font-bold text-charcoal">
              Edit Profile Information
            </h1>
          </div>
          <div className="flex items-center gap-2">
            <button
              type="submit"
              data-slot="save"
              disabled={!isDirty || isSaving || !!nameError || !!gpaError}
              className="rounded-lg bg-dark-green px-4 py-2 font-host text-xs leading-[normal] text-off-white outline-none transition-colors hover:bg-dark-green/90 focus-visible:ring-[3px] focus-visible:ring-dark-green/30 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {isSaving ? "Saving..." : "Save Changes"}
            </button>
          </div>
        </div>

        {(error || status) && (
          <p
            role="status"
            className={cn(
              "font-host text-sm sm:pl-15",
              error ? "text-danger" : "text-dark-green",
            )}
          >
            {error ? errorMessage(error) : status}
          </p>
        )}

        <div className="flex w-full flex-col gap-6 sm:pl-15">
          <Panel title="Personal Information">
            <Field label="First Name" error={nameError}>
              {(id) => (
                <input
                  id={id}
                  data-slot="profile-input"
                  className={fieldClass}
                  value={values.firstName}
                  onChange={(e) => set("firstName")(e.target.value)}
                  autoComplete="given-name"
                  aria-invalid={!!nameError}
                />
              )}
            </Field>
            <Field label="Last Name">
              {(id) => (
                <input
                  id={id}
                  data-slot="profile-input"
                  className={fieldClass}
                  value={values.lastName}
                  onChange={(e) => set("lastName")(e.target.value)}
                  autoComplete="family-name"
                />
              )}
            </Field>
            <Field label="Email">
              {(id) => (
                <input
                  id={id}
                  data-slot="profile-input"
                  className={fieldClass}
                  value={user.email}
                  readOnly
                  disabled
                  title="Email can't be changed here."
                />
              )}
            </Field>
          </Panel>

          {student && (
            <>
              <Panel title="Academic Overview">
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
                      value={values.year}
                      onChange={set("year")}
                      placeholder="Select year"
                      options={graduationOptions().map((o) => ({
                        value: o.year,
                        label: o.label,
                      }))}
                    />
                  )}
                </Field>
                <Field label="GPA (Unweighted)" error={gpaError}>
                  {(id) => (
                    <input
                      id={id}
                      data-slot="profile-input"
                      className={fieldClass}
                      inputMode="numeric"
                      value={values.gpa}
                      onChange={(e) => set("gpa")(e.target.value)}
                      aria-invalid={!!gpaError}
                    />
                  )}
                </Field>
              </Panel>

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
                <Field label="Campus Setting Preference">
                  {(id) => (
                    // No backend field yet, so this can't be saved.
                    <SelectField
                      id={id}
                      value=""
                      onChange={() => {}}
                      placeholder="Coming soon"
                      disabled
                      options={CAMPUS_OPTIONS.map((c) => ({
                        value: c,
                        label: c,
                      }))}
                    />
                  )}
                </Field>
              </Panel>
            </>
          )}
        </div>
      </form>
    </PageShell>
  );
}

export default function EditProfile() {
  const currentUser = useGetCurrentUser();
  const studentID = currentUser.data?.student_id ?? "";
  const student = useGetStudent(studentID, {
    query: { enabled: studentID !== "" },
  });

  const waitingForStudent = studentID !== "" && !student.data && !student.error;
  if (currentUser.isLoading || waitingForStudent) {
    return (
      <PageShell>
        <p className="font-host text-sm">Loading...</p>
      </PageShell>
    );
  }

  const loadError = currentUser.error ?? student.error;
  if (loadError || !currentUser.data) {
    return (
      <PageShell>
        <p className="font-host text-sm text-danger">
          Could not load your profile: {errorMessage(loadError)}
        </p>
      </PageShell>
    );
  }

  return <ProfileForm user={currentUser.data} student={student.data} />;
}
