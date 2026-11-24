import { useState } from "react";
import { Button } from "@/components/ui/button";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { edit } from "../ExtracurricularShared";

// Runs a store call and clears the form only if it succeeded; failures are
// surfaced through the store's `error`.
export const runThen = (action: Promise<unknown>, onDone: () => void) =>
  action.then(onDone, () => undefined);

export const byNewest = <T,>(items: T[], key: (item: T) => string) =>
  [...items].sort((a, b) => key(b).localeCompare(key(a)));

// Idle: [Discard Progress] [Add …]. Editing an item: [Delete …] [Save Edits],
// where Delete asks for confirmation first (Figma: "Delete Note").
export function FormActions({
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
