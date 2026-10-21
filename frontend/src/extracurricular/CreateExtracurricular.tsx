import { useQueryClient } from "@tanstack/react-query";
import { Link, useNavigate } from "react-router-dom";
import {
  getListExtracurricularsQueryKey,
  useCreateExtracurricular,
} from "../api/endpoints";
import ExtracurricularForm from "./ExtracurricularForm";
import {
  TEST_STUDENT_ID,
  emptyFormValues,
  toRequestBody,
  type ExtracurricularFormValues,
} from "./ExtracurricularShared";
import "./Extracurriculars.css";

export default function CreateExtracurricular() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { mutate, isPending, error } = useCreateExtracurricular();

  const handleSubmit = (values: ExtracurricularFormValues) => {
    mutate(
      { studentID: TEST_STUDENT_ID, data: toRequestBody(values) },
      {
        onSuccess: (created) => {
          queryClient.invalidateQueries({
            queryKey: getListExtracurricularsQueryKey(),
          });
          navigate(`/extracurriculars/${created.id}`);
        },
      },
    );
  };

  return (
    <div className="ec-page ec-box">
      <div className="ec-header">
        <Link to="/extracurriculars" className="ec-back" aria-label="Back to extracurriculars">
          ←
        </Link>
        <h2 className="ec-title">Add Extracurricular</h2>
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
