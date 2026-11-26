import { useState } from "react";
import { ChevronDown, Pencil, StickyNote } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { cn } from "cn";
import { edit, longDateLabel } from "../ExtracurricularShared";
import { byNewest, FormActions, runThen } from "./FormActions";
import type { NoteItem, NoteStore } from "./itemSources";
import { SectionCard } from "./SectionCard";
import { ErrorMessage } from "@/components/ui/error-message";

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
        <ErrorMessage error={store.error} prefix="Couldn't save note" />
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
