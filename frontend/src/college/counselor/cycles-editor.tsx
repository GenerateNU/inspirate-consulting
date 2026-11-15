import {
  CollegeApplicationCycleCycle as Cycle,
  type CollegeApplicationCycle,
} from "../../api/models";
import calendarIcon from "@/assets/college/calendar.svg";
import calendarDisabledIcon from "@/assets/college/calendar-disabled.svg";
import checkIcon from "@/assets/college/check.svg";
import { cn } from "@/lib/utils";
import { cycleLabel, toDateInput, toISODateTime } from "../utils/shared";

// row order on the counselor's School Details page
export const editorCycleOrder = [
  Cycle.ED,
  Cycle.ED2,
  Cycle.EA,
  Cycle.REA,
  Cycle.RD,
  Cycle.ROLLING,
] as const;

// one row of the editor; dates are YYYY-MM-DD from date inputs
export interface CycleDraft {
  enabled: boolean;
  deadline: string;
  rollingStart: string;
  rollingEnd: string;
}

export type CyclesDraft = Record<Cycle, CycleDraft>;

// editor state from the college's saved cycles (or an empty form)
export const draftFromCycles = (
  cycles: CollegeApplicationCycle[] | null | undefined,
): CyclesDraft =>
  Object.fromEntries(
    editorCycleOrder.map((cycle) => {
      const saved = cycles?.find((c) => c.cycle === cycle);
      return [
        cycle,
        {
          enabled: !!saved,
          deadline: toDateInput(saved?.deadline),
          rollingStart: saved?.rolling_start ?? "",
          rollingEnd: saved?.rolling_end ?? "",
        },
      ];
    }),
  ) as CyclesDraft;

// the first problem the backend would reject, or null when the draft can be saved
export const cyclesDraftError = (draft: CyclesDraft): string | null => {
  for (const cycle of editorCycleOrder) {
    const row = draft[cycle];
    if (!row.enabled) continue;
    if (cycle === Cycle.ROLLING) {
      if (
        row.rollingStart &&
        row.rollingEnd &&
        row.rollingStart > row.rollingEnd
      )
        return "Rolling Admissions has to open before it closes.";
    } else if (!row.deadline) {
      return `Set a deadline for ${cycleLabel(cycle)}.`;
    }
  }
  return null;
};

// the cycles payload for create/update: only the turned-on rows
export const cyclesFromDraft = (
  draft: CyclesDraft,
): CollegeApplicationCycle[] =>
  editorCycleOrder.flatMap((cycle): CollegeApplicationCycle[] => {
    const row = draft[cycle];
    if (!row.enabled) return [];
    if (cycle === Cycle.ROLLING) {
      return [
        {
          cycle,
          rolling_start: row.rollingStart || undefined,
          rolling_end: row.rollingEnd || undefined,
        },
      ];
    }
    return [{ cycle, deadline: toISODateTime(row.deadline) }];
  });

// checkbox from the hi-fi: 20px box, green with a white check when on
function Checkbox({
  checked,
  label,
  disabled,
  onChange,
}: {
  checked: boolean;
  label: string;
  disabled: boolean;
  onChange: (checked: boolean) => void;
}) {
  return (
    <label className="flex w-[280px] shrink-0 cursor-pointer items-start gap-2 max-md:w-auto">
      <span className="relative flex size-7 shrink-0 items-center justify-center">
        <input
          type="checkbox"
          checked={checked}
          disabled={disabled}
          onChange={(e) => onChange(e.target.checked)}
          className="peer absolute inset-0 cursor-pointer opacity-0"
        />
        <span
          aria-hidden="true"
          className={cn(
            "flex size-5 items-center justify-center rounded border-[1.5px] peer-focus-visible:ring-2 peer-focus-visible:ring-ring",
            checked ? "border-mint bg-mint" : "border-greyish-brown bg-white",
          )}
        >
          {checked && <img src={checkIcon} alt="" width={14} height={14} />}
        </span>
      </span>
      <span className="pt-1 font-host text-sm leading-5 text-charcoal">
        {label}
      </span>
    </label>
  );
}

// date field with the hi-fi calendar icon; the native picker opens from the icon
function DateInput({
  label,
  value,
  enabled,
  disabled,
  widthClassName,
  onChange,
}: {
  label: string;
  value: string;
  enabled: boolean;
  disabled: boolean;
  widthClassName: string;
  onChange: (value: string) => void;
}) {
  return (
    <span className={cn("relative block", widthClassName)}>
      <input
        type="date"
        aria-label={label}
        value={value}
        disabled={!enabled || disabled}
        onChange={(e) => onChange(e.target.value)}
        data-slot="date-input"
        className={cn(
          "relative h-[46px] w-full rounded-lg py-3 pr-12 pl-4 font-host text-sm leading-5 outline-none focus-visible:ring-2 focus-visible:ring-ring",
          "[&::-webkit-calendar-picker-indicator]:absolute [&::-webkit-calendar-picker-indicator]:right-4 [&::-webkit-calendar-picker-indicator]:size-6 [&::-webkit-calendar-picker-indicator]:cursor-pointer [&::-webkit-calendar-picker-indicator]:opacity-0",
          enabled
            ? "border border-greyish-brown bg-white text-charcoal"
            : "cursor-not-allowed border-[1.5px] border-greyish-brown bg-off-white text-greyish-brown",
        )}
      />
      <span className="pointer-events-none absolute top-1/2 right-4 size-6 -translate-y-1/2">
        <img
          src={enabled ? calendarIcon : calendarDisabledIcon}
          alt=""
          width={20}
          height={21}
          className="absolute top-px left-0.5 max-w-none"
        />
      </span>
    </span>
  );
}

interface CyclesEditorProps {
  draft: CyclesDraft;
  disabled?: boolean;
  onChange: (draft: CyclesDraft) => void;
}

// "Application cycles & deadlines": a checkbox per cycle and its deadline
// (or the Opens/Closes window for rolling admissions)
export default function CyclesEditor({
  draft,
  disabled = false,
  onChange,
}: CyclesEditorProps) {
  const update = (cycle: Cycle, change: Partial<CycleDraft>) =>
    onChange({ ...draft, [cycle]: { ...draft[cycle], ...change } });

  return (
    <div className="flex w-full flex-col">
      <div className="flex gap-6 pb-2 font-host text-xs leading-[normal] text-dark-greyish-brown">
        <span className="w-[280px] shrink-0 max-md:w-auto">Cycle</span>
        <span className="flex-1">Deadline</span>
      </div>
      {editorCycleOrder.map((cycle) => {
        const row = draft[cycle];
        const label = cycleLabel(cycle);
        return (
          <div
            key={cycle}
            className="flex items-center gap-6 border-t border-cream py-3 max-md:flex-col max-md:items-start max-md:gap-2"
          >
            <Checkbox
              label={label}
              checked={row.enabled}
              disabled={disabled}
              onChange={(enabled) => update(cycle, { enabled })}
            />
            {cycle === Cycle.ROLLING ? (
              <div className="flex flex-1 flex-wrap gap-4">
                {(
                  [
                    ["Opens", "rollingStart"],
                    ["Closes", "rollingEnd"],
                  ] as const
                ).map(([fieldLabel, key]) => (
                  <span key={key} className="flex flex-col gap-2">
                    <span
                      className={cn(
                        "font-host text-sm leading-5",
                        row.enabled ? "text-charcoal" : "text-greyish-brown",
                      )}
                    >
                      {fieldLabel}
                    </span>
                    <DateInput
                      label={`${label} ${fieldLabel.toLowerCase()}`}
                      value={row[key]}
                      enabled={row.enabled}
                      disabled={disabled}
                      widthClassName="w-[200px]"
                      onChange={(value) => update(cycle, { [key]: value })}
                    />
                  </span>
                ))}
              </div>
            ) : (
              <div className="flex flex-1">
                <DateInput
                  label={`${label} deadline`}
                  value={row.deadline}
                  enabled={row.enabled}
                  disabled={disabled}
                  widthClassName="w-[260px]"
                  onChange={(deadline) => update(cycle, { deadline })}
                />
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
