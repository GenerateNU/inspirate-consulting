import { useQueryClient } from "@tanstack/react-query";
import {
  getGetExtracurricularQueryKey,
  useCreateExtracurricularAchievement,
  useCreateExtracurricularNote,
  useDeleteExtracurricularAchievement,
  useDeleteExtracurricularNote,
  useUpdateExtracurricularAchievement,
  useUpdateExtracurricularNote,
} from "../../api/endpoints";
import type {
  CreateExtracurricularAchievementRequest,
  ExtracurricularAchievement,
  ExtracurricularNote,
  ExtracurricularNoteRequest,
} from "../../api/models";

// Where a section's items live: the API on the edit page, local drafts on
// the create page (posted once the activity exists). `save` rejects on failure
// so the section keeps the user's input.
export type SectionItemSource<TItem, TData> = {
  items: TItem[];
  save: (data: TData, id: number | null) => Promise<unknown>;
  remove: (id: number) => Promise<unknown>;
  busy: boolean;
  error: Error | null;
};

export type AchievementItem = Pick<
  ExtracurricularAchievement,
  "id" | "body" | "entry_date"
>;
export type NoteItem = Pick<ExtracurricularNote, "id" | "body" | "created_at">;
export type AchievementStore = SectionItemSource<
  AchievementItem,
  CreateExtracurricularAchievementRequest
>;
export type NoteStore = SectionItemSource<NoteItem, ExtracurricularNoteRequest>;

function useInvalidateDetail(id: number) {
  const queryClient = useQueryClient();
  return () =>
    queryClient.invalidateQueries({
      queryKey: getGetExtracurricularQueryKey(id),
    });
}

export function useServerAchievements(
  extracurricularId: number,
  achievements: ExtracurricularAchievement[] | null,
): AchievementStore {
  const invalidate = useInvalidateDetail(extracurricularId);
  const options = { mutation: { onSuccess: invalidate } };
  const create = useCreateExtracurricularAchievement(options);
  const update = useUpdateExtracurricularAchievement(options);
  const remove = useDeleteExtracurricularAchievement(options);
  return {
    items: achievements ?? [],
    save: (data, achievementId) =>
      achievementId === null
        ? create.mutateAsync({ id: extracurricularId, data })
        : update.mutateAsync({ id: extracurricularId, achievementId, data }),
    remove: (achievementId) =>
      remove.mutateAsync({ id: extracurricularId, achievementId }),
    busy: create.isPending || update.isPending || remove.isPending,
    error: create.error ?? update.error ?? remove.error,
  };
}

export function useServerNotes(
  extracurricularId: number,
  notes: ExtracurricularNote[] | null,
): NoteStore {
  const invalidate = useInvalidateDetail(extracurricularId);
  const options = { mutation: { onSuccess: invalidate } };
  const create = useCreateExtracurricularNote(options);
  const update = useUpdateExtracurricularNote(options);
  const remove = useDeleteExtracurricularNote(options);
  return {
    items: notes ?? [],
    save: (data, noteId) =>
      noteId === null
        ? create.mutateAsync({ id: extracurricularId, data })
        : update.mutateAsync({ id: extracurricularId, noteId, data }),
    remove: (noteId) => remove.mutateAsync({ id: extracurricularId, noteId }),
    busy: create.isPending || update.isPending || remove.isPending,
    error: create.error ?? update.error ?? remove.error,
  };
}
