import { Link, useParams } from "react-router-dom";
import { useGetExtracurricular } from "../api/endpoints";
import { activityPageClass } from "./ActivityForm";
import ActivityView from "./ActivityView";
import { ErrorMessage } from "@/components/ui/error-message";

export default function ViewExtracurricular() {
  const { id } = useParams();
  const numericId = Number(id);
  const { data, isLoading, error, refetch } = useGetExtracurricular(numericId, {
    query: { enabled: Number.isInteger(numericId) },
  });

  if (isLoading) return <p className={activityPageClass}>Loading...</p>;
  if (error || !data) {
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
            <Link to="/extracurriculars" className="underline">
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
      backTo="/extracurriculars"
      backLabel="Extracurriculars"
      editTo={`/extracurriculars/${data.id}/edit`}
    />
  );
}
