import { useState } from "react";

import {
  CreatePersonalCollegeApplicationRequestBodyApplicationType as ApplicationType,
  CreatePersonalCollegeApplicationRequestBodyCategory as Category,
  type GlobalCollege,
} from "../../api/models";
import { inputClassName, primaryButtonClassName } from "../shared";
import {
  applicationTypes,
  categories,
  categoryLabel,
  typeLabel,
} from "./application-labels";

// pick a college, deadline, and category to add to My List
// the page passes only colleges not already on the list and handles errors
export default function AddApplicationForm({
  colleges,
  initialCollegeId,
  busy,
  adding,
  onAdd,
}: {
  colleges: GlobalCollege[];
  initialCollegeId: number | null;
  busy: boolean;
  adding: boolean;
  onAdd: (
    globalCollegeId: number,
    applicationType: ApplicationType,
    category: Category,
  ) => void;
}) {
  const [collegeId, setCollegeId] = useState<number | null>(initialCollegeId);
  const [applicationType, setApplicationType] = useState<ApplicationType | "">("",);
  const [category, setCategory] = useState<Category | "">("");

  const handleSubmit = (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (collegeId == null || !applicationType || !category) return;
    onAdd(collegeId, applicationType, category);
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="mt-4 flex flex-col gap-3 border border-black p-4"
    >
      <div className="flex flex-wrap items-center gap-2">
        <select
          aria-label="College"
          value={collegeId ?? ""}
          onChange={(e) => setCollegeId(Number(e.target.value) || null)}
          required
          className={inputClassName}
        >
          <option value="">Select a college</option>
          {colleges.map((college) => (
            <option key={college.id} value={college.id}>
              {college.school_name} ({college.school_location})
            </option>
          ))}
        </select>
        <select
          aria-label="Application type"
          value={applicationType}
          onChange={(e) =>
            setApplicationType(e.target.value as ApplicationType)
          }
          required
          className={inputClassName}
        >
          <option value="">Deadline</option>
          {applicationTypes.map((type) => (
            <option key={type} value={type}>
              {typeLabel(type)}
            </option>
          ))}
        </select>
        <select
          aria-label="Category"
          value={category}
          onChange={(e) => setCategory(e.target.value as Category)}
          required
          className={inputClassName}
        >
          <option value="">Category</option>
          {categories.map((category) => (
            <option key={category} value={category}>
              {categoryLabel(category)}
            </option>
          ))}
        </select>
        <button type="submit" disabled={busy} className={primaryButtonClassName}>
          {adding ? "Adding..." : "Add"}
        </button>
      </div>
    </form>
  );
}
