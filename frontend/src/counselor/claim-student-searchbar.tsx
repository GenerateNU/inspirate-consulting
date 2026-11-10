import { ChevronDown } from "lucide-react";
import { SearchBar } from "../../components/Searchbar";

export type Grade = "freshman" | "sophomore" | "junior" | "senior";

const GRADES: { value: Grade; label: string }[] = [
  { value: "freshman", label: "Freshman" },
  { value: "sophomore", label: "Sophomore" },
  { value: "junior", label: "Junior" },
  { value: "senior", label: "Senior" },
];

export function ClaimStudentSearchbar({
  search,
  onSearchChange,
  grade,
  onGradeChange,
}: {
  search: string;
  onSearchChange: (text: string) => void;
  grade: Grade | "";
  onGradeChange: (grade: Grade | "") => void;
}) {
  return (
    <div className="flex items-center gap-4 border border-[#c6c6cd] bg-white p-4">
      <SearchBar
        value={search}
        onChangeText={onSearchChange}
        placeholder="Filter by student name"
        className="flex-1"
      />
      <div className="relative w-[235px] shrink-0">
        <select
          aria-label="Filter by grade"
          value={grade}
          onChange={(e) => onGradeChange(e.target.value as Grade | "")}
          className="h-11 w-full appearance-none rounded-none border border-[#76777d] bg-white py-0 pr-[33px] pl-[13px] font-inter text-sm leading-5 font-medium tracking-[0.14px] text-[#171c20]"
        >
          <option value="">All Grades</option>
          {GRADES.map(({ value, label }) => (
            <option key={value} value={value}>
              {label}
            </option>
          ))}
        </select>
        <ChevronDown
          size={14}
          aria-hidden="true"
          className="pointer-events-none absolute top-1/2 right-[10px] -translate-y-1/2 text-[#171c20]"
        />
      </div>
    </div>
  );
}
