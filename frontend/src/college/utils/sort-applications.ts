import {
  CreatePersonalCollegeApplicationRequestBodyCategory as Category,
  UpdatePersonalCollegeApplicationRequestBodyInterest as Interest,
  type GlobalCollege,
  type PersonalCollegeApplication,
} from "../../api/models";
import { cycleCloseDate, findCycle } from "./shared";

export const sortOptions = [
  { value: "deadline", label: "Deadline" },
  { value: "name", label: "Name" },
  { value: "category", label: "Classification" },
  { value: "interest", label: "Interest" },
] as const;

export type SortKey = (typeof sortOptions)[number]["value"];

const categoryOrder: string[] = [
  Category.reach,
  Category.target,
  Category.safety,
];
const interestOrder: string[] = [
  Interest.top_choice,
  Interest.strong_contender,
  Interest.considering,
];

// position in a fixed order; values not in it (e.g. no interest) go last
const orderOf = (order: string[], value: string | null) => {
  const index = value == null ? -1 : order.indexOf(value);
  return index === -1 ? order.length : index;
};

// a student's applications in the chosen order; ties go by school name
export function sortApplications(
  applications: PersonalCollegeApplication[],
  sortKey: SortKey,
  collegeById: Map<number, GlobalCollege>,
): PersonalCollegeApplication[] {
  const name = (app: PersonalCollegeApplication) =>
    collegeById.get(app.global_college_id)?.school_name ?? "";
  const byName = (
    a: PersonalCollegeApplication,
    b: PersonalCollegeApplication,
  ) => name(a).localeCompare(name(b));

  // closing date of the cycle the student picked; schools without one go last
  const closeDate = (app: PersonalCollegeApplication) =>
    cycleCloseDate(
      findCycle(collegeById.get(app.global_college_id), app.application_type),
    ) ?? Number.POSITIVE_INFINITY;

  const compare: Record<
    SortKey,
    (a: PersonalCollegeApplication, b: PersonalCollegeApplication) => number
  > = {
    name: byName,
    deadline: (a, b) => closeDate(a) - closeDate(b),
    category: (a, b) =>
      orderOf(categoryOrder, a.category) - orderOf(categoryOrder, b.category),
    interest: (a, b) =>
      orderOf(interestOrder, a.interest) - orderOf(interestOrder, b.interest),
  };

  return [...applications].sort(
    (a, b) => compare[sortKey](a, b) || byName(a, b),
  );
}
