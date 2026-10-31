import { CreatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType } from "../../api/models";

// single list of deadline types: the create and update application enums share these values
export const deadlineOptions = [
  { value: ApplicationType.EA, key: "ea_deadline", label: "Early Action" },
  { value: ApplicationType.ED, key: "ed_deadline", label: "Early Decision" },
  { value: ApplicationType.RD, key: "rd_deadline", label: "Regular Decision" },
] as const;

// only need date of full datetimes
export const formatDate = (dateTime: string) =>
  new Date(dateTime).toLocaleDateString(undefined, {
    month: "short",
    day: "numeric",
    year: "numeric",
    timeZone: "UTC",
  });

// the backend wants full datetimes, date inputs give YYYY-MM-DD
export const toISODateTime = (dateStr: string): string | undefined =>
  dateStr ? `${dateStr}T00:00:00Z` : undefined;

// date inputs need YYYY-MM-DD
export const toDateInput = (dateTime?: string | null): string =>
  dateTime ? dateTime.slice(0, 10) : "";

export const inputClassName =
  "min-w-0 flex-1 rounded-md border border-gray-400 px-3 py-2 text-sm disabled:cursor-not-allowed disabled:bg-gray-100";

export const primaryButtonClassName =
  "rounded-md bg-black px-5 py-2 text-sm text-white hover:bg-gray-800 disabled:cursor-not-allowed disabled:opacity-40";

// gray text button, e.g. Edit, Unrank, and the rank arrows on My List
export const smallButtonClassName =
  "text-xs text-gray-500 hover:text-black hover:underline disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:no-underline";

// round outline shared by CircleButton and CircleBadge
export const circleClassName =
  "flex h-8 w-8 shrink-0 items-center justify-center rounded-full border border-black";
