import type {
  CreateExtracurricularRequest,
  Extracurricular,
} from "../api/models";

export const TEST_STUDENT_ID = "00000000-0000-0000-0000-000000000002";

export const CATEGORY_LABELS: Record<
  CreateExtracurricularRequest["type"],
  string
> = {
  maintenance: "Maintenance",
  investment: "Investment",
};

export const STATUS_LABELS: Record<
  CreateExtracurricularRequest["status"],
  string
> = {
  doing: "Currently doing",
  have_done: "Completed",
};

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

export function toRequestBody(
  values: ExtracurricularFormValues,
): CreateExtracurricularRequest {
  return {
    ...values,
    leadership_role: values.leadership_role || undefined,
    end_date: values.end_date || undefined,
  };
}

export function cardSubtitle(item: Extracurricular) {
  const category =
    CATEGORY_LABELS[item.type as ExtracurricularFormValues["type"]] ??
    item.type;
  return [item.leadership_role, item.organization, category]
    .filter(Boolean)
    .join(", ");
}
