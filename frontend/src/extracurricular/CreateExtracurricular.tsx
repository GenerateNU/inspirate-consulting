import { useQueryClient } from "@tanstack/react-query";
import { Link, useNavigate } from "react-router-dom";
import {
  getListExtracurricularsQueryKey,
  useCreateExtracurricular,
} from "../api/endpoints";
import ExtracurricularForm from "./ExtracurricularForm";
import {
  emptyFormValues,
  toRequestBody,
  type ExtracurricularFormValues,
  ec,
} from "./ExtracurricularShared";

export default function CreateExtracurricular() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { mutate, isPending, error } = useCreateExtracurricular();

  const handleSubmit = (values: ExtracurricularFormValues) => {
    // The backend gets the student from the signed-in user, so only the body is sent.
    mutate(
      { data: toRequestBody(values) },
      {
        onSuccess: (created) => {
          void queryClient.invalidateQueries({
            queryKey: getListExtracurricularsQueryKey(),
          });
          void navigate(`/extracurriculars/${created.id}`);
        },
      },
    );
  };

  return (
    <div className={`${ec.page} ${ec.box}`}>
      <div className={ec.header}>
        <Link
          to="/extracurriculars"
          className={ec.back}
          aria-label="Back to extracurriculars"
        >
          ←
        </Link>
        <h2 className={ec.title}>Add Extracurricular</h2>
      </div>
      <ExtracurricularForm
        initialValues={emptyFormValues}
        submitLabel="Save"
        pendingLabel="Saving..."
        isPending={isPending}
        error={error}
        onSubmit={handleSubmit}
      />
    </div>
  );
}
