import { useEffect, useRef, useState } from "react";
import { ChevronDown, X } from "lucide-react";

import type { CounselorStudent } from "../api/models";
import { graduationYearOptions } from "@/lib/graduation-year";
import { cn } from "@/lib/utils";
import { tagClassName } from "../college/utils/components/tag";
import { toggleListItem } from "./toggle-list-items";
import { checkboxClassName, fieldLabelClassName } from "./task-form-styles";

export interface AssignToValue {
  all: boolean;
  graduationYears: number[];
  studentIds: string[];
}

interface AssignToSectionProps {
  value: AssignToValue;
  onChange: (value: AssignToValue) => void;
  students: CounselorStudent[];
  studentsLoading: boolean;
  // editing an individual task: who it's assigned to can't change
  readOnly?: boolean;
}

// Students the task goes to: everyone, the picked classes, plus the picked students
export function resolveStudentIds(
  value: AssignToValue,
  students: CounselorStudent[],
): string[] {
  if (value.all) return students.map((s) => s.student_id);
  const ids = new Set(value.studentIds);
  for (const student of students) {
    if (value.graduationYears.includes(student.graduation_year))
      ids.add(student.student_id);
  }
  return [...ids];
}

// A single student picked by name (no classes) is an individual task
export const isIndividual = (value: AssignToValue) =>
  !value.all &&
  value.graduationYears.length === 0 &&
  value.studentIds.length === 1;

const firstName = (name: string) => name.trim().split(/\s+/)[0] ?? name;

function StudentSearch({
  value,
  onChange,
  students,
  studentsLoading,
  readOnly,
}: AssignToSectionProps) {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const containerRef = useRef<HTMLDivElement>(null);

  // close the list when clicking anywhere else
  useEffect(() => {
    if (!open) return;
    const close = (event: MouseEvent) => {
      if (!containerRef.current?.contains(event.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", close);
    return () => document.removeEventListener("mousedown", close);
  }, [open]);

  const coveredByClass = (student: CounselorStudent) =>
    value.all || value.graduationYears.includes(student.graduation_year);

  const search = query.trim().toLowerCase();
  const matches = students.filter((s) => s.name.toLowerCase().includes(search));
  const picked = students.filter((s) =>
    value.studentIds.includes(s.student_id),
  );
  const selectedCount = resolveStudentIds(value, students).length;

  const toggle = (id: string) =>
    onChange({ ...value, studentIds: toggleListItem(value.studentIds, id) });

  return (
    <div ref={containerRef} className="flex flex-col gap-2">
      <div className="relative">
        <input
          type="text"
          data-slot="student-search"
          role="combobox"
          aria-expanded={open}
          aria-controls="student-search-list"
          aria-label="Search students"
          placeholder="Search students"
          value={query}
          disabled={readOnly}
          onFocus={() => setOpen(true)}
          onChange={(e) => {
            setQuery(e.target.value);
            setOpen(true);
          }}
          className={cn(
            "h-10 w-full rounded-md border bg-white py-2 pr-10 pl-3 font-host text-xs text-charcoal outline-none placeholder:text-light-black disabled:cursor-not-allowed disabled:opacity-60",
            open
              ? "border-light-black bg-off-white"
              : "border-greyish-brown focus-visible:border-light-black",
          )}
        />
        <button
          type="button"
          data-slot="icon-button"
          aria-label={open ? "Close student list" : "Open student list"}
          disabled={readOnly}
          onClick={() => setOpen((o) => !o)}
          className="absolute top-1/2 right-3 -translate-y-1/2 text-light-black disabled:opacity-60"
        >
          <ChevronDown
            className={cn("size-4 transition-transform", open && "rotate-180")}
          />
        </button>
      </div>

      {open && (
        <div
          id="student-search-list"
          className="flex flex-col rounded-md border border-cream bg-white"
        >
          <ul
            role="listbox"
            className="flex max-h-36 flex-col overflow-y-auto p-2"
          >
            {studentsLoading && (
              <li className="px-3 py-2 font-host text-xs text-warm-grey-dark">
                Loading students...
              </li>
            )}
            {!studentsLoading && matches.length === 0 && (
              <li className="px-3 py-2 font-host text-xs text-warm-grey-dark">
                No students match &ldquo;{query}&rdquo;.
              </li>
            )}
            {matches.map((student) => {
              const covered = coveredByClass(student);
              return (
                <li key={student.student_id}>
                  <label
                    className={cn(
                      "flex cursor-pointer items-center gap-2.5 rounded-md px-3 py-2 font-host text-xs text-charcoal hover:bg-off-white",
                      covered && "cursor-default",
                    )}
                  >
                    <input
                      type="checkbox"
                      className={checkboxClassName}
                      checked={
                        covered || value.studentIds.includes(student.student_id)
                      }
                      disabled={covered}
                      onChange={() => toggle(student.student_id)}
                    />
                    {student.name}
                    {covered && (
                      <span className="text-warm-grey">
                        {value.all
                          ? "(all students)"
                          : `(Class of ${student.graduation_year})`}
                      </span>
                    )}
                  </label>
                </li>
              );
            })}
          </ul>
          <p className="border-t border-cream px-4 py-2 font-host text-[10px] text-warm-grey-dark">
            {selectedCount} selected
          </p>
        </div>
      )}

      {!open && picked.length > 0 && (
        <ul aria-label="Selected students" className="flex flex-wrap gap-2">
          {picked.map((student) => (
            <li
              key={student.student_id}
              className={cn(
                tagClassName,
                "gap-1.5 bg-lightest-green text-[10px] text-charcoal",
              )}
            >
              {firstName(student.name)}
              {!readOnly && (
                <button
                  type="button"
                  data-slot="icon-button"
                  aria-label={`Remove ${student.name}`}
                  onClick={() => toggle(student.student_id)}
                  className="text-charcoal hover:text-dark-green"
                >
                  <X className="size-3" strokeWidth={3} />
                </button>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

// "Assign to": whole classes, specific students, or both
export default function AssignToSection(props: AssignToSectionProps) {
  const { value, onChange, readOnly } = props;
  const classes = graduationYearOptions();

  return (
    <div className="flex flex-col gap-6">
      <fieldset className="flex flex-col gap-3">
        <legend className={cn(fieldLabelClassName, "mb-3")}>
          Class groups
        </legend>
        <div className="flex flex-wrap items-center gap-x-6 gap-y-2">
          {classes.map(({ value: year }) => (
            <label
              key={year}
              className="flex items-center gap-2 font-host text-xs text-charcoal"
            >
              <input
                type="checkbox"
                className={checkboxClassName}
                checked={value.all || value.graduationYears.includes(year)}
                disabled={readOnly || value.all}
                onChange={() =>
                  onChange({
                    ...value,
                    graduationYears: value.graduationYears.includes(year)
                      ? value.graduationYears.filter((y) => y !== year)
                      : [...value.graduationYears, year],
                  })
                }
              />
              Class of {year}
            </label>
          ))}
          <label className="flex items-center gap-2 font-host text-xs text-charcoal">
            <input
              type="checkbox"
              className={checkboxClassName}
              checked={value.all}
              disabled={readOnly}
              onChange={(e) =>
                onChange({
                  ...value,
                  all: e.target.checked,
                  graduationYears: [],
                })
              }
            />
            All
          </label>
        </div>
      </fieldset>

      <div className="flex flex-col gap-3">
        <span className={fieldLabelClassName}>Individual students</span>
        <StudentSearch {...props} />
      </div>
    </div>
  );
}
