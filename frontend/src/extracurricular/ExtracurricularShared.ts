import type {
  CreateExtracurricularRequest,
  CreateExtracurricularRequestStatus as Status,
  CreateExtracurricularRequestType as Category,
  Extracurricular,
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

// Form fields are all strings so inputs stay controlled; optional fields use "".
export type ExtracurricularFormValues = {
  name: string;
  leadership_role: string;
  organization: string;
  type: Category;
  status: Status;
  start_date: string;
  end_date: string;
  description: string;
};

export const emptyFormValues: ExtracurricularFormValues = {
  name: "",
  leadership_role: "",
  organization: "",
  type: "maintenance",
  status: "doing",
  start_date: "",
  end_date: "",
  description: "",
};

// The API returns full ISO datetimes; date inputs need yyyy-MM-dd.
export const toDateInput = (value?: string | null) =>
  value ? value.slice(0, 10) : "";

export function toFormValues(item: Extracurricular): ExtracurricularFormValues {
  return {
    name: item.name,
    leadership_role: item.leadership_role ?? "",
    organization: item.organization,
    type: isCategory(item.type) ? item.type : emptyFormValues.type,
    status: isStatus(item.status) ? item.status : emptyFormValues.status,
    start_date: toDateInput(item.start_date),
    end_date: toDateInput(item.end_date),
    description: item.description,
  };
}

// Optional fields are left out of the request when empty.
export function toRequestBody(
  values: ExtracurricularFormValues,
): CreateExtracurricularRequest {
  return {
    ...values,
    leadership_role: values.leadership_role || undefined,
    end_date: values.end_date || undefined,
  };
}

// Matches the lo-fi card subtitle: "Position, Organization, Category".
export function cardSubtitle(item: Extracurricular) {
  return [item.leadership_role, item.organization, categoryLabel(item.type)]
    .filter(Boolean)
    .join(", ");
}

// Tailwind class strings shared across the extracurricular pages.
export const ec = {
  page: "mx-auto max-w-5xl px-6 py-8",
  box: "border border-current px-8 py-6",
  header: "mb-6 flex items-center gap-3",
  title: "text-3xl",
  back: "text-2xl",
  linkButton: "inline-block border border-current px-3 py-1.5 text-xs",
  list: "grid gap-4",
  card: "flex items-center justify-between gap-4 border border-current px-5 py-3.5",
  cardText: "flex flex-col gap-1",
  cardTitle: "font-semibold",
  cardSubtitle: "text-xs",
  section: "mt-5 border border-current px-5 py-4",
  sectionHeading: "mb-3 text-base",
  formGrid: "grid gap-x-8 gap-y-4 md:grid-cols-3",
  field: "flex flex-col gap-1 text-sm",
};
