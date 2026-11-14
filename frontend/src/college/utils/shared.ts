import {
  CollegeApplicationCycleCycle as Cycle,
  type CollegeApplicationCycle,
  type GlobalCollege,
} from "../../api/models";

// every application cycle in the order the backend returns them
export const cycleOptions = [
  { value: Cycle.ED, label: "Early Decision" },
  { value: Cycle.ED2, label: "Early Decision II" },
  { value: Cycle.REA, label: "Restrictive Early Action" },
  { value: Cycle.EA, label: "Early Action" },
  { value: Cycle.RD, label: "Regular Decision" },
  { value: Cycle.ROLLING, label: "Rolling Admissions" },
] as const;

export const cycleLabel = (cycle: string) =>
  cycleOptions.find(({ value }) => value === cycle)?.label ?? cycle;

// the college's row for one cycle, or undefined when it doesn't offer it
export const findCycle = (
  college: GlobalCollege | undefined,
  cycle: string,
): CollegeApplicationCycle | undefined =>
  college?.cycles?.find((offered) => offered.cycle === cycle);

// only need date of full datetimes
export const formatDate = (dateTime: string) =>
  new Date(dateTime).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
    timeZone: "UTC",
  });

// "Nov 1, 2026" for a deadline, "Sep 1, 2026 – Mar 1, 2027" for a rolling window
export const cycleDateLabel = (cycle: CollegeApplicationCycle): string => {
  if (cycle.deadline) return formatDate(cycle.deadline);
  const { rolling_start: start, rolling_end: end } = cycle;
  if (start && end) return `${formatDate(start)} – ${formatDate(end)}`;
  if (start) return `Opens ${formatDate(start)}`;
  if (end) return `Until ${formatDate(end)}`;
  return "";
};

// the date a cycle closes, for sorting; rolling uses the end of its window
export const cycleCloseDate = (
  cycle: CollegeApplicationCycle | undefined,
): number | undefined => {
  const date = cycle?.deadline ?? cycle?.rolling_end;
  return date ? new Date(date).getTime() : undefined;
};

// the backend wants full datetimes, date inputs give YYYY-MM-DD
export const toISODateTime = (dateStr: string): string | undefined =>
  dateStr ? `${dateStr}T00:00:00Z` : undefined;

// date inputs need YYYY-MM-DD
export const toDateInput = (dateTime?: string | null): string =>
  dateTime ? dateTime.slice(0, 10) : "";
