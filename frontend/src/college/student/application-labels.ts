import {
  CreatePersonalCollegeApplicationRequestBodyCategory as Category,
  PersonalCollegeApplicationEssayProgress as EssayProgress,
  UpdatePersonalCollegeApplicationRequestBodyInterest as Interest,
} from "../../api/models";
import type { ChoiceOption } from "../utils/components/choice-chips";
import type { TagTone } from "../utils/components/tag";

// order matches the School Details page
export const categoryOptions: ChoiceOption<Category>[] = [
  { value: Category.target, label: "Target" },
  { value: Category.reach, label: "Reach" },
  { value: Category.safety, label: "Safety" },
];

export const interestOptions: ChoiceOption<Interest>[] = [
  { value: Interest.top_choice, label: "Top Choice" },
  { value: Interest.strong_contender, label: "Strong Contender" },
  { value: Interest.considering, label: "Considering" },
];

export const categoryLabel = (category: string) =>
  categoryOptions.find(({ value }) => value === category)?.label ?? category;

export const interestLabel = (interest: string) =>
  interestOptions.find(({ value }) => value === interest)?.label ?? interest;

// the pill next to a school's name on My List
export const essayProgressTag: Record<
  string,
  { label: string; tone: TagTone }
> = {
  [EssayProgress.not_started]: { label: "Essays Not Started", tone: "grey" },
  [EssayProgress.in_progress]: { label: "Essays In Progress", tone: "yellow" },
  [EssayProgress.completed]: { label: "All Essays Completed", tone: "green" },
};
