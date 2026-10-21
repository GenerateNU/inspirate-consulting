import { useQueryClient } from "@tanstack/react-query";
import { Link, useNavigate, useParams } from "react-router-dom";
import { errorMessage } from "../api/client";
import {
  getListExtracurricularsQueryKey,
  useListExtracurriculars,
  useUpdateExtracurricular,
} from "../api/endpoints";
import ExtracurricularForm from "./ExtracurricularForm";
import {
  toFormValues,
  toRequestBody,
  type ExtracurricularFormValues,
} from "./ExtracurricularShared";
import "./Extracurriculars.css";

export default function EditExtracurricular() {
  const { id } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  // There is no "get one" endpoint yet, so find the entry in the student's list.
  const { data, isLoading, error: loadError } = useListExtracurriculars();
  const { mutate, isPending, error } = useUpdateExtracurricular();

  const item = (data ?? []).find((entry) => String(entry.id) === id);

  if (isLoading) return <p className="ec-page">Loading...</p>;
  if (loadError) return <p className="ec-page">Error: {errorMessage(loadError)}</p>;
  if (!item) {
    return (
      <div className="ec-page">
        <p>This extracurricular doesn't exist or was removed.</p>
        <Link to="/extracurriculars">Back to extracurriculars</Link>
      </div>
    );
  }

  const handleSubmit = (values: ExtracurricularFormValues) => {
    mutate(
      { id: item.id, data: toRequestBody(values) },
      {
        onSuccess: () => {
          queryClient.invalidateQueries({
            queryKey: getListExtracurricularsQueryKey(),
          });
          navigate(`/extracurriculars/${item.id}`);
        },
      },
    );
  };

  return (
    <div className="ec-page ec-box">
      <div className="ec-header">
        <Link
          to={`/extracurriculars/${item.id}`}
          className="ec-back"
          aria-label={`Back to ${item.name}`}
        >
          ←
        </Link>
        <h2 className="ec-title">Edit {item.name}</h2>
      </div>
      <ExtracurricularForm
        key={item.id}
        initialValues={toFormValues(item)}
        submitLabel="Save changes"
        pendingLabel="Saving..."
        isPending={isPending}
        error={error}
        onSubmit={handleSubmit}
      />
    </div>
  );
}
