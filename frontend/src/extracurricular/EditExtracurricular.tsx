import { useQueryClient } from "@tanstack/react-query";
import { Link, useLocation, useNavigate, useParams } from "react-router-dom";
import { errorMessage } from "../api/client";
import {
  getGetExtracurricularQueryKey,
  getListExtracurricularsQueryKey,
  useGetExtracurricular,
  useUpdateExtracurricular,
} from "../api/endpoints";
import type { ExtracurricularDetail } from "../api/models";
import ActivityForm, { activityPageClass } from "./ActivityForm";
import type { CreateNotice } from "./CreateExtracurricular";
import {
  AchievementsSection,
  NotesSection,
  useServerAchievements,
  useServerNotes,
} from "./EditSections";
import {
  toActivityFormValues,
  toUpdateRequestBody,
  type ActivityFormValues,
} from "./ExtracurricularShared";

export default function EditExtracurricular() {
  const { id } = useParams();
  const numericId = Number(id);
  const { data, isLoading, error } = useGetExtracurricular(numericId, {
    query: { enabled: Number.isInteger(numericId) },
  });

  if (isLoading) return <p className={activityPageClass}>Loading...</p>;
  if (error || !data) {
    return (
      <div className={activityPageClass}>
        <p>
          {error && error.status !== 404
            ? `Error: ${errorMessage(error)}`
            : "This extracurricular doesn't exist or was removed."}
        </p>
        <Link to="/extracurriculars" className="underline">
          Back to extracurriculars
        </Link>
      </div>
    );
  }

  return <EditPage key={data.id} item={data} />;
}

function EditPage({ item }: { item: ExtracurricularDetail }) {
  const navigate = useNavigate();
  const location = useLocation();
  const notice = (location.state as CreateNotice | null)?.notice;
  const queryClient = useQueryClient();
  const { mutate, isPending, error } = useUpdateExtracurricular();
  const achievements = useServerAchievements(item.id, item.achievements);
  const notes = useServerNotes(item.id, item.notes);

  const handleSubmit = (values: ActivityFormValues) => {
    mutate(
      { id: item.id, data: toUpdateRequestBody(values) },
      {
        onSuccess: async () => {
          await Promise.all([
            queryClient.invalidateQueries({
              queryKey: getListExtracurricularsQueryKey(),
            }),
            queryClient.invalidateQueries({
              queryKey: getGetExtracurricularQueryKey(item.id),
            }),
          ]);
          void navigate(`/extracurriculars/${item.id}`);
        },
      },
    );
  };

  return (
    <ActivityForm
      title={item.name}
      showTitleIcon
      notice={notice}
      breadcrumb={{ label: item.name, to: `/extracurriculars/${item.id}` }}
      initialValues={toActivityFormValues(item)}
      isPending={isPending}
      error={error}
      onSubmit={handleSubmit}
      achievements={<AchievementsSection store={achievements} />}
      notes={<NotesSection store={notes} />}
    />
  );
}
