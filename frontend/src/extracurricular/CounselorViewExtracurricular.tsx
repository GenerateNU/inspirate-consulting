import { Link, useParams } from "react-router-dom";
import { useGetExtracurricular } from "../api/endpoints";
import { useStudentProfile } from "../counselor/hooks/use-student-profile";
import { activityPageClass } from "./ActivityForm";
import ActivityView from "./ActivityView";
import { ErrorMessage } from "@/components/ui/error-message";

// a counselor's read-only view of one of their student's extracurriculars
export default function CounselorViewExtracurricular() {
  const { studentId = "", id } = useParams();
  const numericId = Number(id);
  const { name } = useStudentProfile(studentId);
  const { data, isLoading, error, refetch } = useGetExtracurricular(numericId, {
    query: { enabled: Number.isInteger(numericId) },
  });
  const listPath = `/counselor/students/${studentId}/extracurriculars`;

  if (isLoading) return <p className={activityPageClass}>Loading...</p>;
  // the activity is fetched by id alone, so also check it's this student's
  if (error || !data || data.student_id !== studentId) {
    const loadFailed = !!error && error.status !== 404;
    return (
      <div className={activityPageClass}>
        <ErrorMessage
          variant="banner"
          title={loadFailed ? "Couldn't load this extracurricular" : undefined}
          error={error}
          message={
            loadFailed
              ? undefined
              : "This extracurricular doesn't exist or was removed."
          }
          onRetry={loadFailed ? refetch : undefined}
          action={
            <Link to={listPath} className="underline">
              Back to extracurriculars
            </Link>
          }
        />
      </div>
    );
  }

  return (
    <ActivityView
      item={data}
      backTo={listPath}
      backLabel={name ? `${name}'s Extracurriculars` : "Extracurriculars"}
    />
  );
}
