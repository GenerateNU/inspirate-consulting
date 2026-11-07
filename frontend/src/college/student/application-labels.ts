import {
  CreatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType,
  CreatePersonalCollegeApplicationRequestBodyCategory as Category,
} from "../../api/models";
import { deadlineOptions } from "../utils/shared";

export const applicationTypes = Object.values(ApplicationType);
export const categories = Object.values(Category);

export const typeLabel = (type: string) =>
  deadlineOptions.find(({ value }) => value === type)?.label ?? type;

export const categoryLabel = (category: string) =>
  category.charAt(0).toUpperCase() + category.slice(1);

export const deadlineClassName = "border-blue-300 bg-blue-100 text-blue-700";

export const categoryClassName: Record<string, string> = {
  [Category.target]: "border-red-300 bg-red-100 text-red-700",
  [Category.reach]: "border-amber-300 bg-amber-100 text-amber-700",
  [Category.safety]: "border-emerald-300 bg-emerald-100 text-emerald-700",
};
