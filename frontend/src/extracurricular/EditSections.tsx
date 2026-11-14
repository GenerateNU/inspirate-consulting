import { useState, type ReactNode } from "react";
import { useQueryClient } from "@tanstack/react-query";
import type { LucideIcon } from "lucide-react";
import { ChevronDown, Pencil, StickyNote, Trophy } from "lucide-react";
import { errorMessage } from "../api/client";
import {
  getGetExtracurricularQueryKey,
  useCreateExtracurricularAchievement,
  useCreateExtracurricularNote,
  useDeleteExtracurricularAchievement,
  useDeleteExtracurricularNote,
  useUpdateExtracurricularAchievement,
  useUpdateExtracurricularNote,
} from "../api/endpoints";
import type {
  CreateExtracurricularAchievementRequest,
  ExtracurricularAchievement,
  ExtracurricularNote,
  ExtracurricularNoteRequest,
} from "../api/models";
import { Button } from "@/components/ui/button";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { cn } from "cn";
import { edit, longDateLabel, toDateInput } from "./ExtracurricularShared";

export function SectionCard({
  icon: Icon,
  title,
  className,
  children,
}: {
  icon: LucideIcon;
  title: string;
  className?: string;
  children: ReactNode;
}) {
  return (
    <section className={cn(edit.card, className)}>
      <div className="flex items-center gap-5">
        <CategoryIcon icon={Icon} />
        <h2 className={edit.cardTitle}>{title}</h2>
      </div>
      {children}
    </section>
  );
}

export function CategoryIcon({ icon: Icon }: { icon: LucideIcon }) {
  return (
    <span
      aria-hidden="true"
      className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-[#f8efe4] text-charcoal"
    >
      <Icon className="size-6" strokeWidth={1.75} />
    </span>
  );
}

// Where a section's items live: the API on the edit page, local drafts on
// the create page (posted once the activity exists). `save` rejects on failure
// so the section keeps the user's input.
export type ItemStore<TItem, TData> = {
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
export type AchievementStore = ItemStore<
  AchievementItem,
  CreateExtracurricularAchievementRequest
>;
export type NoteStore = ItemStore<NoteItem, ExtracurricularNoteRequest>;

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

// Runs a store call and clears the form only if it succeeded; failures are
// surfaced through the store's `error`.
const runThen = (action: Promise<unknown>, onDone: () => void) =>
  action.then(onDone, () => undefined);

// Idle: [Discard Progress] [Add …]. Editing an item: [Delete …] [Save Edits],
// where Delete asks for confirmation first (Figma: "Delete Note").
function FormActions({
  itemLabel,
  isEditing,
  onDiscard,
  onConfirmDelete,
  busy,
}: {
  itemLabel: string;
  isEditing: boolean;
  onDiscard?: () => void;
  onConfirmDelete: () => Promise<unknown>;
  busy: boolean;
}) {
  const [confirming, setConfirming] = useState(false);

  return (
    <div className="flex flex-wrap items-center justify-end gap-3">
      {isEditing ? (
        <Button
          type="button"
          className={edit.dangerButton}
          onClick={() => setConfirming(true)}
          disabled={busy}
        >
          Delete {itemLabel}
        </Button>
      ) : (
        onDiscard && (
          <Button
            type="button"
            variant="outline"
            className={edit.secondaryButton}
            onClick={onDiscard}
            disabled={busy}
          >
            Discard Progress
          </Button>
        )
      )}
      <Button type="submit" className={edit.primaryButton} disabled={busy}>
        {isEditing ? "Save Edits" : `Add ${itemLabel}`}
      </Button>
      <ConfirmDialog
        open={confirming}
        onOpenChange={setConfirming}
        title={`Are you sure you want to delete this ${itemLabel.toLowerCase()}?`}
        description="Once deleted, this cannot be recovered."
        confirmLabel={`Delete ${itemLabel}`}
        pendingLabel="Deleting..."
        pending={busy}
        onConfirm={() =>
          void onConfirmDelete().finally(() => setConfirming(false))
        }
      />
    </div>
  );
}

const byNewest = <T,>(items: T[], key: (item: T) => string) =>
  [...items].sort((a, b) => key(b).localeCompare(key(a)));

export function AchievementsSection({
  store,
  showDiscardWhenIdle = true,
}: {
  store: AchievementStore;
  // The create design only offers "Discard" while editing an achievement.
  showDiscardWhenIdle?: boolean;
}) {
  const [editingId, setEditingId] = useState<number | null>(null);
  const [entryDate, setEntryDate] = useState("");
  const [body, setBody] = useState("");

  const reset = () => {
    setEditingId(null);
    setEntryDate("");
    setBody("");
  };

  // The pencil toggles: clicking it again on the same item stops editing,
  // since the editing state has no Discard button.
  const startEditing = (achievement: AchievementItem) => {
    if (editingId === achievement.id) return reset();
    setEditingId(achievement.id);
    setEntryDate(toDateInput(achievement.entry_date));
    setBody(achievement.body);
  };

  const handleSubmit = (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    void runThen(
      store.save({ body: body.trim(), entry_date: entryDate }, editingId),
      reset,
    );
  };

  const sorted = byNewest(store.items, (a) => a.entry_date);

  return (
    <SectionCard icon={Trophy} title="Achievements & Awards">
      {sorted.length > 0 && (
        <ul className="flex flex-wrap gap-4">
          {sorted.map((achievement) => (
            <AchievementTile
              key={achievement.id}
              achievement={achievement}
              isEditing={editingId === achievement.id}
              onEdit={() => startEditing(achievement)}
            />
          ))}
        </ul>
      )}

      <form
        onSubmit={handleSubmit}
        className="flex flex-col gap-3 rounded-lg border-[0.5px] border-[#bdb7ad] bg-white p-6"
      >
        <h3 className="text-base leading-5 text-charcoal">
          {editingId === null
            ? "Add a New Achievement or Award"
            : "Edit Your Achievement or Award"}
        </h3>
        <div className="flex flex-col gap-2">
          <Label htmlFor="ec-achievement-date" className={edit.label}>
            Date
          </Label>
          <Input
            id="ec-achievement-date"
            type="date"
            className={cn(edit.input, !entryDate && edit.emptyDate)}
            value={entryDate}
            onChange={(e) => setEntryDate(e.target.value)}
            required
          />
        </div>
        <div className="flex flex-col gap-2">
          <Label htmlFor="ec-achievement-body" className={edit.label}>
            Description
          </Label>
          <Input
            id="ec-achievement-body"
            className={edit.input}
            placeholder="Achievement Description"
            value={body}
            onChange={(e) => setBody(e.target.value)}
            required
          />
        </div>
        {store.error && (
          <p className={edit.error}>Error: {errorMessage(store.error)}</p>
        )}
        <FormActions
          itemLabel="Achievement"
          isEditing={editingId !== null}
          onDiscard={showDiscardWhenIdle ? reset : undefined}
          onConfirmDelete={() =>
            editingId === null
              ? Promise.resolve()
              : runThen(store.remove(editingId), reset)
          }
          busy={store.busy}
        />
      </form>
    </SectionCard>
  );
}

export function NotesSection({
  store,
  className,
}: {
  store: NoteStore;
  className?: string;
}) {
  const [editingId, setEditingId] = useState<number | null>(null);
  const [body, setBody] = useState("");

  const reset = () => {
    setEditingId(null);
    setBody("");
  };

  const handleSubmit = (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    void runThen(store.save({ body: body.trim() }, editingId), reset);
  };

  const sorted = byNewest(store.items, (n) => n.created_at);

  return (
    <SectionCard
      icon={StickyNote}
      title="Notes & Scratch Work"
      className={className}
    >
      {sorted.length > 0 && (
        <ul className="flex flex-col gap-3">
          {sorted.map((note) => (
            <NoteRow
              key={note.id}
              note={note}
              isEditing={editingId === note.id}
              onEdit={() => {
                if (editingId === note.id) return reset();
                setEditingId(note.id);
                setBody(note.body);
              }}
            />
          ))}
        </ul>
      )}

      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <div className="flex flex-col gap-2">
          <Label htmlFor="ec-note-body" className={edit.label}>
            {editingId === null
              ? "Write some notes you think are important:"
              : "Edit Your Note"}
          </Label>
          <Textarea
            id="ec-note-body"
            className={cn(edit.input, "min-h-[84px]")}
            placeholder="Notes"
            value={body}
            onChange={(e) => setBody(e.target.value)}
            required
          />
        </div>
        {store.error && (
          <p className={edit.error}>Error: {errorMessage(store.error)}</p>
        )}
        <FormActions
          itemLabel="Note"
          isEditing={editingId !== null}
          onDiscard={reset}
          onConfirmDelete={() =>
            editingId === null
              ? Promise.resolve()
              : runThen(store.remove(editingId), reset)
          }
          busy={store.busy}
        />
      </form>
    </SectionCard>
  );
}

export function AchievementTile({
  achievement,
  isEditing = false,
  onEdit,
}: {
  achievement: AchievementItem;
  isEditing?: boolean;
  onEdit?: () => void;
}) {
  return (
    <li
      className={cn(
        "flex w-40 flex-col items-start gap-3 rounded-lg border border-[#bdb7ad] bg-white p-3",
        isEditing && "border-[#5e5c56]",
      )}
    >
      <p className="w-full text-sm leading-5 break-words text-charcoal">
        {achievement.body}
      </p>
      <span className={edit.tag}>
        {longDateLabel(achievement.entry_date, { calendarDate: true })}
      </span>
      {onEdit && (
        <Button
          type="button"
          variant="ghost"
          size="icon"
          className={cn(edit.iconButton, "self-end")}
          onClick={onEdit}
          aria-label={`Edit achievement: ${achievement.body}`}
          aria-pressed={isEditing}
        >
          <Pencil aria-hidden="true" />
        </Button>
      )}
    </li>
  );
}

export function Toggle({
  checked,
  onChange,
  labelledBy,
}: {
  checked: boolean;
  // Leave out to render a read-only switch.
  onChange?: (checked: boolean) => void;
  labelledBy: string;
}) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      aria-labelledby={labelledBy}
      aria-readonly={onChange ? undefined : true}
      onClick={onChange ? () => onChange(!checked) : undefined}
      className={cn(
        "relative h-8 w-[52px] shrink-0 rounded-full transition-colors outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50",
        checked ? "bg-[#19664c]" : "bg-charcoal",
        !onChange && "cursor-default",
      )}
    >
      <span
        aria-hidden="true"
        className={cn(
          "absolute top-1 left-1 size-6 rounded-full bg-white transition-transform",
          checked && "translate-x-5",
        )}
      />
    </button>
  );
}

// Notes longer than this collapse to one line with an expand chevron.
const NOTE_PREVIEW_CHARS = 60;

export function NoteRow({
  note,
  isEditing = false,
  onEdit,
}: {
  note: NoteItem;
  isEditing?: boolean;
  onEdit?: () => void;
}) {
  const [expanded, setExpanded] = useState(false);
  const collapsible =
    note.body.length > NOTE_PREVIEW_CHARS || note.body.includes("\n");

  return (
    <li
      className={cn(
        "flex items-start gap-2 rounded-lg border border-[#bdb7ad] bg-white py-2 pr-2 pl-1",
        isEditing && "border-[#5e5c56]",
      )}
    >
      <span
        aria-hidden="true"
        className="w-[3px] self-stretch rounded-xl bg-tan"
      />
      <div className="min-w-0 flex-1 text-sm leading-5 text-charcoal">
        <p>{longDateLabel(note.created_at)}</p>
        <p
          className={
            expanded || !collapsible
              ? "break-words whitespace-pre-wrap"
              : "truncate"
          }
        >
          {note.body}
        </p>
      </div>
      <div className="flex flex-col gap-1">
        {onEdit && (
          <Button
            type="button"
            variant="ghost"
            size="icon"
            className={edit.iconButton}
            onClick={onEdit}
            aria-label="Edit note"
            aria-pressed={isEditing}
          >
            <Pencil aria-hidden="true" />
          </Button>
        )}
        {collapsible && (
          <Button
            type="button"
            variant="ghost"
            size="icon"
            className={edit.iconButton}
            onClick={() => setExpanded((prev) => !prev)}
            aria-expanded={expanded}
            aria-label={expanded ? "Collapse note" : "Expand note"}
          >
            <ChevronDown
              aria-hidden="true"
              className={cn("transition-transform", expanded && "rotate-180")}
            />
          </Button>
        )}
      </div>
    </li>
  );
}
