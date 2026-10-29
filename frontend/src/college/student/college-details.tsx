import CollegeDetailsLayout, {
  useCollegeFromRoute,
} from "../utils/college-details-layout";
import { Field, deadlineOptions, formatDate } from "../utils/shared";

// read-only view of one college
export default function StudentCollegeDetails() {
  const { data: college, isLoading, error } = useCollegeFromRoute();

  return (
    <CollegeDetailsLayout college={college} isLoading={isLoading} error={error}>
      {college && (
        <dl>
          <Field label="Location">{college.school_location}</Field>
          {deadlineOptions.map(({ value, key, label }) => (
            <Field key={value} label={`${label} (${value})`}>
              {college[key] ? formatDate(college[key]!) : "N/A"}
            </Field>
          ))}
        </dl>
      )}
    </CollegeDetailsLayout>
  );
}
