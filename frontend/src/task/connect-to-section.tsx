import { useEffect, useRef, useState, type ReactNode } from "react";
import { ChevronDown, ChevronRight } from "lucide-react";

import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { cn } from "@/lib/utils";
import {
  fieldLabelClassName,
  hintClassName,
  radioClassName,
} from "./task-form-styles";
import {
  useTaskLinkOptions,
  type EssayGroup,
  type LinkOption,
  type TaskLinkType,
} from "./use-task-link-options";

export interface ConnectToValue {
  type: TaskLinkType | null;
  // the specific item, or "" for none
  itemId: string;
}

interface ConnectToSectionProps {
  value: ConnectToValue;
  onChange: (value: ConnectToValue) => void;
  // everyone the task is going to
  studentIds: string[];
  // true for one student picked by name
  individual: boolean;
}

const TYPE_OPTIONS: { value: TaskLinkType; label: string }[] = [
  { value: "college", label: "Colleges" },
  { value: "essay", label: "Essays" },
  { value: "extracurricular", label: "Extracurriculars" },
  { value: "module", label: "Modules" },
];

// essays and extracurriculars belong to one student, so the backend only bulk-links colleges and modules
export const INDIVIDUAL_ONLY_TYPES: TaskLinkType[] = [
  "essay",
  "extracurricular",
];

const PICKER_LABELS: Partial<Record<TaskLinkType, string>> = {
  college: "Individual college (optional)",
  essay: "Individual essay (optional)",
  module: "Module (optional)",
};

function OptionRadio({
  name,
  option,
  checked,
  onSelect,
}: {
  name: string;
  option: LinkOption;
  checked: boolean;
  onSelect: (id: string) => void;
}) {
  return (
    <label
      className={cn(
        "flex cursor-pointer items-center gap-2.5 rounded-md px-3 py-2 font-host text-xs text-charcoal hover:bg-off-white",
        checked && "bg-lightest-green hover:bg-lightest-green",
      )}
    >
      <input
        type="radio"
        name={name}
        className={radioClassName}
        checked={checked}
        // clicking the picked item again unlinks it, since the item is optional
        onClick={() => onSelect(checked ? "" : option.id)}
        onChange={() => {}}
      />
      {option.label}
    </label>
  );
}

function EssayGroups({
  groups,
  selectedId,
  onSelect,
}: {
  groups: EssayGroup[];
  selectedId: string;
  onSelect: (id: string) => void;
}) {
  return (
    <>
      {groups.map((group) => (
        <Collapsible
          key={group.collegeId ?? "other"}
          // open the college holding the picked essay, or the first one
          defaultOpen={
            group.essays.some((essay) => essay.id === selectedId) ||
            (!selectedId && group === groups[0])
          }
        >
          <CollapsibleTrigger asChild>
            <button
              type="button"
              data-slot="essay-group-trigger"
              className="group flex w-full items-center gap-2 rounded-md bg-off-white px-3 py-2 text-left font-host text-[10px] text-charcoal"
            >
              <ChevronRight className="size-3 transition-transform group-data-[state=open]:rotate-90" />
              <span className="flex-1">{group.collegeName}</span>
              <span className="text-warm-grey-dark">
                {group.essays.length}{" "}
                {group.essays.length === 1 ? "essay" : "essays"}
              </span>
            </button>
          </CollapsibleTrigger>
          <CollapsibleContent className="flex flex-col pl-4">
            {group.essays.map((essay) => (
              <OptionRadio
                key={essay.id}
                name="task-essay"
                option={essay}
                checked={essay.id === selectedId}
                onSelect={onSelect}
              />
            ))}
          </CollapsibleContent>
        </Collapsible>
      ))}
    </>
  );
}

// A dropdown-looking button that opens a list under it
function ItemPicker({
  label,
  selectedLabel,
  placeholder,
  children,
}: {
  label: string;
  selectedLabel: string | undefined;
  placeholder: string;
  children: ReactNode;
}) {
  const [open, setOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const close = (event: MouseEvent) => {
      if (!containerRef.current?.contains(event.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", close);
    return () => document.removeEventListener("mousedown", close);
  }, [open]);

  return (
    <div ref={containerRef} className="flex flex-col gap-3">
      <span className={fieldLabelClassName}>{label}</span>
      <button
        type="button"
        data-slot="item-picker-trigger"
        aria-expanded={open}
        onClick={() => setOpen((o) => !o)}
        className={cn(
          "flex h-10 w-full items-center justify-between gap-2 rounded-md border bg-white px-3 text-left font-host text-xs",
          open ? "border-light-black" : "border-greyish-brown",
        )}
      >
        <span
          className={cn(
            "truncate",
            selectedLabel ? "text-charcoal" : "text-warm-grey",
          )}
        >
          {selectedLabel ?? placeholder}
        </span>
        <ChevronDown
          className={cn(
            "size-4 shrink-0 text-light-black transition-transform",
            open && "rotate-180",
          )}
        />
      </button>
      {open && (
        <div className="flex max-h-56 flex-col gap-1 overflow-y-auto rounded-md border border-cream bg-white p-2">
          {children}
        </div>
      )}
    </div>
  );
}

function EmptyNote({ children }: { children: ReactNode }) {
  return (
    <p className="px-3 py-2 font-host text-xs text-warm-grey-dark">
      {children}
    </p>
  );
}

// "Connected to": what the task is about, then (optionally) the specific item
export default function ConnectToSection({
  value,
  onChange,
  studentIds,
  individual,
}: ConnectToSectionProps) {
  const options = useTaskLinkOptions(value.type, studentIds);
  const selectItem = (itemId: string) => onChange({ ...value, itemId });

  const flatOptions =
    value.type === "college"
      ? options.colleges
      : value.type === "module"
        ? options.modules
        : options.essayGroups.flatMap((group) => group.essays);
  const selectedLabel = flatOptions.find((o) => o.id === value.itemId)?.label;

  const pickerLabel = value.type ? PICKER_LABELS[value.type] : undefined;
  const showPicker =
    pickerLabel !== undefined &&
    studentIds.length > 0 &&
    (value.type !== "essay" || individual);

  const renderOptions = () => {
    if (options.isLoading) return <EmptyNote>Loading...</EmptyNote>;

    if (value.type === "essay") {
      return options.essayGroups.length === 0 ? (
        <EmptyNote>This student doesn&apos;t have any essays yet.</EmptyNote>
      ) : (
        <EssayGroups
          groups={options.essayGroups}
          selectedId={value.itemId}
          onSelect={selectItem}
        />
      );
    }

    const list = value.type === "college" ? options.colleges : options.modules;
    if (list.length === 0) {
      return (
        <EmptyNote>
          {value.type === "college"
            ? individual
              ? "This student hasn't added any colleges yet."
              : "No colleges in common."
            : individual
              ? "This student has finished every module."
              : "No modules yet."}
        </EmptyNote>
      );
    }
    return list.map((option) => (
      <OptionRadio
        key={option.id}
        name={`task-${value.type}`}
        option={option}
        checked={option.id === value.itemId}
        onSelect={selectItem}
      />
    ));
  };

  return (
    <div className="flex flex-col gap-6">
      <fieldset className="flex flex-col gap-3">
        <legend className={cn(fieldLabelClassName, "mb-3")}>
          Connected to
        </legend>
        <div className="flex flex-wrap items-center gap-x-6 gap-y-2">
          {TYPE_OPTIONS.map((option) => {
            const disabled =
              !individual && INDIVIDUAL_ONLY_TYPES.includes(option.value);
            return (
              <label
                key={option.value}
                className={cn(
                  "flex items-center gap-2 font-host text-xs text-charcoal",
                  disabled && "text-warm-grey",
                )}
              >
                <input
                  type="radio"
                  name="task-type"
                  className={radioClassName}
                  checked={value.type === option.value}
                  disabled={disabled}
                  // clicking the picked type again clears it, since a task doesn't need a type
                  onClick={() =>
                    onChange(
                      value.type === option.value
                        ? { type: null, itemId: "" }
                        : { type: option.value, itemId: "" },
                    )
                  }
                  onChange={() => {}}
                />
                {option.label}
              </label>
            );
          })}
        </div>
        {!individual && (
          <p className={hintClassName}>
            Essays and extracurriculars belong to one student, so pick a single
            student to connect one.
          </p>
        )}
      </fieldset>

      {value.type === "extracurricular" && individual && (
        <p className={hintClassName}>
          Picking a specific extracurricular isn&apos;t available yet, so this
          task will link to the student&apos;s extracurriculars in general.
        </p>
      )}

      {showPicker && (
        <ItemPicker
          // remount when the type changes so the list starts closed
          key={value.type}
          label={pickerLabel}
          selectedLabel={selectedLabel}
          placeholder={value.type === "essay" ? "Select an essay" : "Select"}
        >
          {renderOptions()}
        </ItemPicker>
      )}
    </div>
  );
}
