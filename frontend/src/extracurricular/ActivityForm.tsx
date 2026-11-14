import { useState, type ReactNode } from "react";
import { Link } from "react-router-dom";
import {
  ArrowLeft,
  ChevronRight,
  Goal,
  Palette,
  Paperclip,
} from "lucide-react";
import { errorMessage } from "../api/client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  NativeSelect,
  NativeSelectOption,
} from "@/components/ui/native-select";
import { Textarea } from "@/components/ui/textarea";
import { cn } from "cn";
import {
  APPLICATION_WORD_LIMIT,
  CATEGORY_LABELS,
  STATUS_LABELS,
  edit,
  isCategory,
  isStatus,
  wordCount,
  type ActivityFormValues,
} from "./ExtracurricularShared";
import { CategoryIcon, SectionCard, Toggle } from "./sections";

const FORM_ID = "ec-activity-form";

export const activityPageClass =
  "flex min-h-full flex-1 flex-col gap-4 bg-[#fbfaf7] px-6 py-12 font-host-grotesk md:px-[114px]";

interface ActivityFormProps {
  title: string;
  // The edit design shows a category icon beside the title; create doesn't.
  showTitleIcon?: boolean;
  submitLabel?: string;
  // One-off message shown above the form, e.g. after a partial create.
  notice?: string;
  // The create design stacks Achievements 8px under Details as one group.
  groupAchievementsWithDetails?: boolean;
  // Last breadcrumb after "Extracurriculars"; links when `to` is set.
  breadcrumb: { label: string; to?: string };
  initialValues: ActivityFormValues;
  isPending: boolean;
  error: Error | null;
  onSubmit: (values: ActivityFormValues) => void;
  // Notes and achievements save on their own, so each page supplies them.
  achievements: ReactNode;
  notes: ReactNode;
}

export default function ActivityForm({
  title,
  showTitleIcon = false,
  submitLabel = "Save",
  notice,
  groupAchievementsWithDetails = false,
  breadcrumb,
  initialValues,
  isPending,
  error,
  onSubmit,
  achievements,
  notes,
}: ActivityFormProps) {
  const [values, setValues] = useState(initialValues);
  const [formError, setFormError] = useState<string | null>(null);

  const update = (patch: Partial<ActivityFormValues>) => {
    setValues((prev) => ({ ...prev, ...patch }));
    setFormError(null);
  };

  const words = wordCount(values.application_description);
  const overWordLimit = words > APPLICATION_WORD_LIMIT;
  const endsBeforeStart =
    values.end_date !== "" &&
    values.start_date !== "" &&
    values.end_date < values.start_date;

  const handleSubmit = (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (endsBeforeStart) {
      setFormError("End date can't be before the start date.");
      return;
    }
    if (overWordLimit) {
      setFormError(
        `Application description must be ${APPLICATION_WORD_LIMIT} words or fewer.`,
      );
      return;
    }
    onSubmit(values);
  };

  return (
    <div className={activityPageClass}>
      <nav
        aria-label="Breadcrumb"
        className="flex items-center gap-2 text-xs text-[#303436]"
      >
        <Link
          to="/extracurriculars"
          className="flex items-center gap-2 hover:underline"
        >
          <ArrowLeft className="size-4" aria-hidden="true" />
          Extracurriculars
        </Link>
        <ChevronRight className="size-4" aria-hidden="true" />
        {breadcrumb.to ? (
          <Link to={breadcrumb.to} className="hover:underline">
            {breadcrumb.label}
          </Link>
        ) : (
          <span aria-current="page">{breadcrumb.label}</span>
        )}
      </nav>

      <div className="flex items-center justify-between gap-4">
        <div className="flex min-w-0 items-center gap-6">
          {showTitleIcon && <CategoryIcon icon={Palette} />}
          <h1 className={cn(edit.pageTitle, "truncate")}>{title}</h1>
        </div>
        <Button
          type="submit"
          form={FORM_ID}
          className={edit.primaryButton}
          disabled={isPending}
        >
          {isPending ? "Saving..." : submitLabel}
        </Button>
      </div>

      {notice && (
        <p
          role="status"
          className="rounded-lg border border-tan bg-[#f8efe4] px-4 py-3 text-sm text-[#786246]"
        >
          {notice}
        </p>
      )}

      {(formError || error) && (
        <p role="alert" className={edit.error}>
          {formError ?? `Error: ${errorMessage(error)}`}
        </p>
      )}

      <div
        className={cn(
          "flex flex-col",
          groupAchievementsWithDetails ? "gap-2" : "gap-4",
        )}
      >
        {/* Inputs in the other cards join this form via the `form` attribute,
          since the notes and achievements forms can't be nested inside it. */}
        <form
          id={FORM_ID}
          onSubmit={handleSubmit}
          className="flex flex-col gap-2"
        >
          <h2 className={edit.cardTitle}>Details</h2>
          <div className="grid gap-x-4 gap-y-6 rounded-xl border border-cream bg-white p-6 md:grid-cols-3">
            <Field id="ec-name" label="Activity / Sport Name">
              <Input
                id="ec-name"
                className={edit.input}
                placeholder="Activity Name"
                value={values.name}
                onChange={(e) => update({ name: e.target.value })}
                required
              />
            </Field>
            <Field id="ec-role" label="Position">
              <Input
                id="ec-role"
                className={edit.input}
                placeholder="Position"
                maxLength={50}
                value={values.leadership_role}
                onChange={(e) => update({ leadership_role: e.target.value })}
              />
            </Field>
            <Field id="ec-org" label="Organization">
              <Input
                id="ec-org"
                className={edit.input}
                placeholder="Organization"
                maxLength={100}
                value={values.organization}
                onChange={(e) => update({ organization: e.target.value })}
                required
              />
            </Field>
            <Field id="ec-hours" label="Hours Per Week">
              <Input
                id="ec-hours"
                type="number"
                placeholder="Eg. 10"
                inputMode="numeric"
                min={0}
                max={168}
                className={edit.input}
                value={values.hours_per_week}
                onChange={(e) => update({ hours_per_week: e.target.value })}
              />
            </Field>
            <Field id="ec-type" label="Category">
              <NativeSelect
                id="ec-type"
                className={cn(edit.input, "pr-10")}
                value={values.type}
                onChange={(e) => {
                  const type = e.target.value;
                  if (isCategory(type)) update({ type });
                }}
              >
                {Object.entries(CATEGORY_LABELS).map(([value, label]) => (
                  <NativeSelectOption key={value} value={value}>
                    {label}
                  </NativeSelectOption>
                ))}
              </NativeSelect>
            </Field>
            <Field id="ec-weeks" label="Weeks Per Year">
              <Input
                id="ec-weeks"
                type="number"
                placeholder="Eg. 10"
                inputMode="numeric"
                min={0}
                max={52}
                className={edit.input}
                value={values.weeks_per_year}
                onChange={(e) => update({ weeks_per_year: e.target.value })}
              />
            </Field>
            <Field id="ec-status" label="Status">
              <NativeSelect
                id="ec-status"
                className={cn(edit.input, "pr-10")}
                value={values.status}
                onChange={(e) => {
                  const status = e.target.value;
                  if (isStatus(status)) update({ status });
                }}
              >
                {Object.entries(STATUS_LABELS).map(([value, label]) => (
                  <NativeSelectOption key={value} value={value}>
                    {label}
                  </NativeSelectOption>
                ))}
              </NativeSelect>
            </Field>
            <Field id="ec-start" label="Start Date">
              <Input
                id="ec-start"
                type="date"
                className={cn(edit.input, !values.start_date && edit.emptyDate)}
                value={values.start_date}
                onChange={(e) => update({ start_date: e.target.value })}
                required
              />
            </Field>
            <Field id="ec-end" label="End Date (optional)">
              <Input
                id="ec-end"
                type="date"
                min={values.start_date || undefined}
                className={cn(edit.input, !values.end_date && edit.emptyDate)}
                value={values.end_date}
                onChange={(e) => update({ end_date: e.target.value })}
                aria-invalid={endsBeforeStart || undefined}
              />
            </Field>
          </div>
        </form>

        {achievements}
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        {notes}
        <SectionCard icon={Goal} title="Goals">
          <div className="flex flex-1 flex-col gap-2">
            <Label htmlFor="ec-goals" className={edit.label}>
              Write Down Your Goal Here
            </Label>
            <Textarea
              id="ec-goals"
              form={FORM_ID}
              className={cn(edit.input, "min-h-[156px] flex-1")}
              placeholder="Goals"
              value={values.goals}
              onChange={(e) => update({ goals: e.target.value })}
            />
          </div>
        </SectionCard>
      </div>

      <SectionCard icon={Paperclip} title="Application Information">
        <div className="flex flex-col gap-2">
          <Label htmlFor="ec-app-description" className={edit.label}>
            Application description - up to {APPLICATION_WORD_LIMIT} words
            (Common App)
          </Label>
          <Textarea
            id="ec-app-description"
            form={FORM_ID}
            className={cn(edit.input, "min-h-[176px]")}
            placeholder="Application Description"
            value={values.application_description}
            onChange={(e) =>
              update({ application_description: e.target.value })
            }
            aria-invalid={overWordLimit || undefined}
            aria-describedby="ec-app-word-count"
          />
          <p
            id="ec-app-word-count"
            className={cn(
              "text-right text-xs text-[#303436]",
              overWordLimit && "text-destructive",
            )}
          >
            {words} / {APPLICATION_WORD_LIMIT} words
          </p>
        </div>
        <hr className="border-cream" />
        <div className="flex items-center justify-between gap-4 px-2 py-1">
          <span
            id="ec-ready-label"
            className="text-base leading-5 text-charcoal"
          >
            Ready For College Application?
          </span>
          <Toggle
            labelledBy="ec-ready-label"
            checked={values.ready_for_application}
            onChange={(ready_for_application) =>
              update({ ready_for_application })
            }
          />
        </div>
      </SectionCard>
    </div>
  );
}

function Field({
  id,
  label,
  children,
}: {
  id: string;
  label: string;
  children: ReactNode;
}) {
  return (
    <div className="flex flex-col gap-2 [&_[data-slot=native-select-wrapper]]:w-full">
      <Label htmlFor={id} className={edit.label}>
        {label}
      </Label>
      {children}
    </div>
  );
}
