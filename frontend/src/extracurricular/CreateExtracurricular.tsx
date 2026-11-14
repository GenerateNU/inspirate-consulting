import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import {
  createExtracurricularAchievement,
  createExtracurricularNote,
  getListExtracurricularsQueryKey,
  useCreateExtracurricular,
} from "../api/endpoints";
import ActivityForm from "./ActivityForm";
import {
  AchievementsSection,
  NotesSection,
  type AchievementItem,
  type ItemStore,
  type NoteItem,
} from "./EditSections";
import {
  emptyActivityValues,
  toCreateRequestBody,
  type ActivityFormValues,
} from "./ExtracurricularShared";

export type CreateNotice = { notice: string };

// Achievements and notes need the activity's id, so on this page they're kept
// as local drafts and posted right after the activity is created.
function useDraftStore<TItem extends { id: number }, TData>(
  toItem: (data: TData, id: number) => TItem,
): ItemStore<TItem, TData> {
  const [items, setItems] = useState<TItem[]>([]);
  const [nextId, setNextId] = useState(1);
  return {
    items,
    save: (data, id) => {
      if (id === null) {
        setItems((prev) => [...prev, toItem(data, nextId)]);
        setNextId((n) => n + 1);
      } else {
        setItems((prev) =>
          prev.map((item) => (item.id === id ? toItem(data, id) : item)),
        );
      }
      return Promise.resolve();
    },
    remove: (id) => {
      setItems((prev) => prev.filter((item) => item.id !== id));
      return Promise.resolve();
    },
    busy: false,
    error: null,
  };
}

export default function CreateExtracurricular() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { mutate, isPending, error } = useCreateExtracurricular();
  const [savingDrafts, setSavingDrafts] = useState(false);

  const achievements = useDraftStore(
    (data: { body: string; entry_date: string }, id): AchievementItem => ({
      id,
      ...data,
    }),
  );
  const notes = useDraftStore((data: { body: string }, id): NoteItem => ({
    id,
    body: data.body,
    created_at: new Date().toISOString(),
  }));

  const handleSubmit = (values: ActivityFormValues) => {
    // The backend gets the student from the signed-in user, so only the body is sent.
    mutate(
      { data: toCreateRequestBody(values) },
      {
        onSuccess: async (created) => {
          setSavingDrafts(true);
          // Oldest first so created_at order matches the order they were added.
          const results = await Promise.allSettled([
            ...achievements.items.map(({ body, entry_date }) =>
              createExtracurricularAchievement(created.id, {
                body,
                entry_date,
              }),
            ),
            ...notes.items.map(({ body }) =>
              createExtracurricularNote(created.id, { body }),
            ),
          ]);
          void queryClient.invalidateQueries({
            queryKey: getListExtracurricularsQueryKey(),
          });
          const failed = results.filter((r) => r.status === "rejected").length;
          if (failed > 0) {
            // The activity exists now, so send them to edit it rather than
            // letting a retry here create a duplicate.
            const state: CreateNotice = {
              notice: `Activity saved, but ${failed} achievement(s) or note(s) didn't save. Please add them again.`,
            };
            void navigate(`/extracurriculars/${created.id}/edit`, {
              replace: true,
              state,
            });
          } else {
            void navigate(`/extracurriculars/${created.id}`, { replace: true });
          }
        },
      },
    );
  };

  return (
    <ActivityForm
      title="Add Extracurricular"
      submitLabel="Add Extracurricular"
      groupAchievementsWithDetails
      breadcrumb={{ label: "Add Extracurricular" }}
      initialValues={emptyActivityValues}
      isPending={isPending || savingDrafts}
      error={error}
      onSubmit={handleSubmit}
      achievements={
        <AchievementsSection store={achievements} showDiscardWhenIdle={false} />
      }
      notes={<NotesSection store={notes} />}
    />
  );
}
