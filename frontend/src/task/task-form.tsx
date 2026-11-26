import { useState, type ReactNode } from "react";
import { useNavigate } from "react-router-dom";

import { useListCounselorStudents } from "../api/endpoints";
import arrowLeftIcon from "@/assets/college/arrow-left.svg";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { cn } from "@/lib/utils";
import AssignToSection, {
  isIndividual,
  resolveStudentIds,
  type AssignToValue,
} from "./assign-to-section";
import ConnectToSection, {
  INDIVIDUAL_ONLY_TYPES,
  type ConnectToValue,
} from "./connect-to-section";
import {
  fieldInputClassName,
  fieldLabelClassName,
  hintClassName,
} from "./task-form-styles";
import { ErrorMessage } from "@/components/ui/error-message";

export const DESCRIPTION_WORD_LIMIT = 250;
const TITLE_MAX_LENGTH = 200;

export const countWords = (text: string) =>
  text.trim() === "" ? 0 : text.trim().split(/\s+/).length;

export interface TaskFormValues {
  title: string;
  description: string;
  // YYYY-MM-DD from the date input, or ""
  dueDate: string;
  assignTo: AssignToValue;
  connectTo: ConnectToValue;
}

interface TaskFormProps {
  heading: string;
  initialValues: TaskFormValues;
  submitLabel: string;
  pendingLabel: string;
  isPending: boolean;
  error: unknown;
  // editing a task: who it's assigned to can't change
  assigneeReadOnly?: boolean;
  // the description can't be cleared once set
  requireDescription?: boolean;
  onSubmit: (values: TaskFormValues, studentIds: string[]) => void;
}

function Section({
  title,
  subtitle,
  children,
}: {
  title: string;
  subtitle?: string;
  children: ReactNode;
}) {
  return (
    <section className="flex flex-col gap-6 border-b border-cream px-6 py-8 last:border-b-0">
      <div className="flex flex-col gap-1">
        <h2 className="font-host text-base leading-[normal] text-charcoal">
          {title}
        </h2>
        {subtitle && <p className={hintClassName}>{subtitle}</p>}
      </div>
      {children}
    </section>
  );
}

// The hi-fi Create/Edit Task page: task details, who it's for, and what it's connected to
export default function TaskForm({
  heading,
  initialValues,
  submitLabel,
  pendingLabel,
  isPending,
  error,
  assigneeReadOnly = false,
  requireDescription = false,
  onSubmit,
}: TaskFormProps) {
  const navigate = useNavigate();
  const [values, setValues] = useState(initialValues);
  const { title, description, dueDate, assignTo, connectTo } = values;

  const studentsQuery = useListCounselorStudents();
  const students = studentsQuery.data ?? [];
  const studentIds = resolveStudentIds(assignTo, students);
  const individual = isIndividual(assignTo);

  const wordCount = countWords(description);
  const overLimit = wordCount > DESCRIPTION_WORD_LIMIT;

  const set = (change: Partial<TaskFormValues>) =>
    setValues((prev) => ({ ...prev, ...change }));

  const changeAssignTo = (next: AssignToValue) => {
    // the picked item belongs to the old students, so clear it;
    // and groups can't link to essays or extracurriculars
    const keepType =
      connectTo.type !== null &&
      (isIndividual(next) || !INDIVIDUAL_ONLY_TYPES.includes(connectTo.type));
    set({
      assignTo: next,
      connectTo: {
        type: keepType ? connectTo.type : null,
        itemId: connectTo.type === "module" ? connectTo.itemId : "",
      },
    });
  };

  const missing = [
    !title.trim() && "a task name",
    requireDescription && !description.trim() && "a description",
    studentIds.length === 0 &&
      !studentsQuery.isLoading &&
      "at least one student",
  ].filter(Boolean);
  const canSubmit = missing.length === 0 && !overLimit && !isPending;

  const handleSubmit = (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (canSubmit) onSubmit(values, studentIds);
  };

  return (
    <div className="flex min-h-full w-full flex-col items-center bg-off-white px-12 py-10 max-lg:px-6">
      <div className="flex w-full max-w-164 flex-col gap-6">
        <div className="relative flex items-center">
          <button
            type="button"
            data-slot="icon-button"
            aria-label="Back"
            onClick={() => navigate(-1)}
            className="absolute -left-16 flex size-8 items-center justify-center rounded-md hover:bg-cream max-lg:static max-lg:mr-3"
          >
            <img src={arrowLeftIcon} alt="" width={24} height={24} />
          </button>
          <h1 className="font-jakarta text-2xl leading-[normal] font-semibold text-charcoal">
            {heading}
          </h1>
        </div>

        <form
          onSubmit={handleSubmit}
          className="flex flex-col rounded-lg border border-cream bg-white shadow-card"
        >
          <Section title="Task details">
            <label className="flex flex-col gap-3">
              <span className={fieldLabelClassName}>Task name</span>
              <Input
                value={title}
                maxLength={TITLE_MAX_LENGTH}
                placeholder={"e.g. Finish “Why Columbia?” draft"}
                onChange={(e) => set({ title: e.target.value })}
                className={cn(fieldInputClassName, "h-10")}
              />
            </label>

            <div className="flex flex-col gap-3">
              <label htmlFor="task-description" className={fieldLabelClassName}>
                Description
              </label>
              <Textarea
                id="task-description"
                value={description}
                placeholder="Write your response..."
                aria-invalid={overLimit}
                aria-describedby="task-description-count"
                onChange={(e) => set({ description: e.target.value })}
                className={cn(fieldInputClassName, "min-h-34")}
              />
              <p
                id="task-description-count"
                className={cn(hintClassName, overLimit && "text-danger")}
              >
                {wordCount} / {DESCRIPTION_WORD_LIMIT} words
              </p>
            </div>

            <label className="flex w-62.5 max-w-full flex-col gap-3">
              <span className={fieldLabelClassName}>Due date</span>
              <Input
                type="date"
                value={dueDate}
                onChange={(e) => set({ dueDate: e.target.value })}
                className={cn(
                  fieldInputClassName,
                  "h-10",
                  !dueDate && "text-warm-grey",
                )}
              />
            </label>
          </Section>

          <Section
            title="Assign to"
            subtitle="Choose a whole class, specific students, or both."
          >
            {studentsQuery.error ? (
              <ErrorMessage
                variant="banner"
                title="Couldn't load your students"
                error={studentsQuery.error}
                onRetry={studentsQuery.refetch}
              />
            ) : (
              <AssignToSection
                value={assignTo}
                onChange={changeAssignTo}
                students={students}
                studentsLoading={studentsQuery.isLoading}
                readOnly={assigneeReadOnly}
              />
            )}
          </Section>

          <Section
            title="Task type"
            subtitle="Pick what this task is about, then link the specific item."
          >
            <ConnectToSection
              value={connectTo}
              onChange={(next) => set({ connectTo: next })}
              studentIds={studentIds}
              individual={individual}
            />

            <div className="flex flex-col items-end gap-2 pt-4">
              <div className="flex gap-2.5">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => navigate(-1)}
                  className="h-7 rounded-md border-greyish-brown px-5 font-host text-xs font-normal text-light-black shadow-none hover:bg-off-white"
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  size="sm"
                  disabled={!canSubmit}
                  className="h-7 rounded-md bg-dark-green px-5 font-host text-xs font-normal text-white hover:bg-forest"
                >
                  {isPending ? pendingLabel : submitLabel}
                </Button>
              </div>
              {missing.length > 0 && (
                <p className={hintClassName}>Add {missing.join(" and ")}.</p>
              )}
              <ErrorMessage error={error} size="xs" />
            </div>
          </Section>
        </form>
      </div>
    </div>
  );
}
