import type {
  CreateExtracurricularRequest,
  Extracurricular,
} from "../api/models";

// TODO: replace with the logged-in student's id from the auth context
export const TEST_STUDENT_ID = "00000000-0000-0000-0000-000000000002";

export const CATEGORY_LABELS: Record<CreateExtracurricularRequest["type"], string> = {
  maintenance: "Maintenance",
  investment: "Investment",
};

export const STATUS_LABELS: Record<CreateExtracurricularRequest["status"], string> = {
  doing: "Currently doing",
  have_done: "Completed",
};

// Form fields are all strings so inputs stay controlled; optional fields use "".
export type ExtracurricularFormValues = {
  name: string;
  leadership_role: string;
  organization: string;
  type: CreateExtracurricularRequest["type"];
  status: CreateExtracurricularRequest["status"];
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
    type: item.type as ExtracurricularFormValues["type"],
    status: item.status as ExtracurricularFormValues["status"],
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
  const category =
    CATEGORY_LABELS[item.type as ExtracurricularFormValues["type"]] ?? item.type;
  return [item.leadership_role, item.organization, category]
    .filter(Boolean)
    .join(", ");
}
