import { useState } from "react";
import { Pencil, Trophy } from "lucide-react";
import { errorMessage } from "../../api/client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "cn";
import { edit, longDateLabel, toDateInput } from "../ExtracurricularShared";
import { byNewest, FormActions, runThen } from "./FormActions";
import type { AchievementItem, AchievementStore } from "./itemSources";
import { SectionCard } from "./SectionCard";

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
