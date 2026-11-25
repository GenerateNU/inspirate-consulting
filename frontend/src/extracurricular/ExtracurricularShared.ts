import type {
  CreateExtracurricularRequest,
  CreateExtracurricularRequestStatus as Status,
  CreateExtracurricularRequestType as Category,
  Extracurricular,
  UpdateExtracurricularRequest,
} from "../api/models";

export type { Category, Status };

export const CATEGORY_LABELS: Record<Category, string> = {
  maintenance: "Maintenance",
  investment: "Investment",
};

export const STATUS_LABELS: Record<Status, string> = {
  doing: "Currently doing",
  have_done: "Completed",
};

export const isCategory = (value: string): value is Category =>
  Object.hasOwn(CATEGORY_LABELS, value);

export const isStatus = (value: string): value is Status =>
  Object.hasOwn(STATUS_LABELS, value);

export const categoryLabel = (value: string) =>
  isCategory(value) ? CATEGORY_LABELS[value] : value;

export const statusLabel = (value: string) =>
  isStatus(value) ? STATUS_LABELS[value] : value;

// The API returns full ISO datetimes; date inputs need yyyy-MM-dd.
export const toDateInput = (value?: string | null) =>
  value ? value.slice(0, 10) : "";

export function dateRangeLabel(item: Extracurricular) {
  const start = item.start_date?.slice(0, 4);
  if (!start) return "";
  if (item.status === "doing" && !item.end_date) return `${start} – Present`;
  const end = item.end_date?.slice(0, 4);
  return end && end !== start ? `${start} – ${end}` : start;
}

export type ActivityFormValues = {
  name: string;
  leadership_role: string;
  organization: string;
  type: Category;
  status: Status;
  start_date: string;
  end_date: string;
  hours_per_week: string;
  weeks_per_year: string;
  goals: string;
  application_description: string;
  ready_for_application: boolean;
};

export const emptyActivityValues: ActivityFormValues = {
  name: "",
  leadership_role: "",
  organization: "",
  type: "investment",
  status: "doing",
  start_date: "",
  end_date: "",
  hours_per_week: "",
  weeks_per_year: "",
  goals: "",
  application_description: "",
  ready_for_application: false,
};

export function toActivityFormValues(
  item: Extracurricular,
): ActivityFormValues {
  return {
    name: item.name,
    leadership_role: item.leadership_role ?? "",
    organization: item.organization,
    type: isCategory(item.type) ? item.type : emptyActivityValues.type,
    status: isStatus(item.status) ? item.status : emptyActivityValues.status,
    start_date: toDateInput(item.start_date),
    end_date: toDateInput(item.end_date),
    hours_per_week: item.hours_per_week?.toString() ?? "",
    weeks_per_year: item.weeks_per_year?.toString() ?? "",
    goals: item.goals ?? "",
    application_description: item.application_description ?? "",
    ready_for_application: item.ready_for_application,
  };
}

const numberOrNull = (value: string) => (value === "" ? null : Number(value));

const numberOrUndefined = (value: string) =>
  value === "" ? undefined : Number(value);

// POST body: empty optional fields are left out.
export function toCreateRequestBody(
  values: ActivityFormValues,
): CreateExtracurricularRequest {
  return {
    ...values,
    leadership_role: values.leadership_role || undefined,
    end_date: values.end_date || undefined,
    hours_per_week: numberOrUndefined(values.hours_per_week),
    weeks_per_year: numberOrUndefined(values.weeks_per_year),
    goals: values.goals || undefined,
    application_description: values.application_description || undefined,
  };
}

// PATCH body: empty optional fields are sent as null so they get cleared.
export function toUpdateRequestBody(
  values: ActivityFormValues,
): UpdateExtracurricularRequest {
  return {
    ...values,
    leadership_role: values.leadership_role || null,
    end_date: values.end_date || null,
    hours_per_week: numberOrNull(values.hours_per_week),
    weeks_per_year: numberOrNull(values.weeks_per_year),
    goals: values.goals || null,
    application_description: values.application_description || null,
  };
}

export const APPLICATION_WORD_LIMIT = 250;

export const wordCount = (text: string) =>
  text.trim() === "" ? 0 : text.trim().split(/\s+/).length;

const ordinal = (day: number) => {
  if (day % 100 >= 11 && day % 100 <= 13) return `${day}th`;
  return `${day}${["th", "st", "nd", "rd"][day % 10] ?? "th"}`;
};

// "November 3rd, 2026". Calendar dates (yyyy-MM-dd...) are read as-is so a
// timezone offset can't shift them a day; full timestamps use local time.
export function longDateLabel(value: string, { calendarDate = false } = {}) {
  const date = calendarDate
    ? new Date(`${value.slice(0, 10)}T00:00:00`)
    : new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  const month = date.toLocaleString("en-US", { month: "long" });
  return `${month} ${ordinal(date.getDate())}, ${date.getFullYear()}`;
}

// Tailwind class strings for the edit page (Figma: Extracurricular - Edit Activity).
export const edit = {
  card: "flex flex-col gap-6 rounded-xl border border-cream bg-white p-6",
  pageTitle: "font-jakarta text-[30px] font-bold text-charcoal",
  cardTitle: "font-jakarta text-xl font-medium text-charcoal",
  label: "text-sm leading-5 font-normal text-charcoal",
  input:
    "h-auto rounded-lg border-[#bdb7ad] bg-white px-4 py-3 text-sm leading-5 text-charcoal shadow-none placeholder:text-[#5e5c56] md:text-sm resize-none [appearance:textfield] [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none",
  // Native date inputs show their mm/dd/yyyy mask as text, not a placeholder.
  emptyDate: "text-[#5e5c56]",
  primaryButton:
    "h-auto rounded-lg bg-[#19664c] px-6 py-2 text-sm leading-5 font-normal text-[#fbfaf7] hover:bg-[#19664c]/90",
  dangerButton:
    "h-auto rounded-lg bg-[#b0132b] px-6 py-2 text-sm leading-5 font-normal text-[#fbfaf7] hover:bg-[#b0132b]/90",
  secondaryButton:
    "h-auto rounded-lg border-[#bdb7ad] bg-[#fbfaf7] px-6 py-2 text-sm leading-5 font-normal text-[#5e5c56] shadow-none",
  iconButton:
    "size-7 rounded-md text-charcoal [&_svg:not([class*='size-'])]:size-5",
  tag: "rounded-full bg-cream px-2 py-1 text-xs leading-none font-medium whitespace-nowrap text-[#5e5c56]",
};
