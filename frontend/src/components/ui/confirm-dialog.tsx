import { useRef, type ReactNode } from "react";
import { AlertDialog } from "radix-ui";
import { XIcon } from "lucide-react";
import { cn } from "cn";

type ConfirmDialogProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: ReactNode;
  description?: ReactNode;
  confirmLabel: string;
  // Shown on the confirm button while `pending`, e.g. "Deleting...".
  pendingLabel?: string;
  cancelLabel?: string;
  // "destructive" (red) for deletes and other irreversible actions.
  variant?: "destructive" | "default";
  // Keeps the dialog open with the confirm button disabled while the
  // action runs. The caller closes it via `onOpenChange(false)` when done.
  pending?: boolean;
  onConfirm: () => void;
};

/**
 * Confirmation popup (Figma: "Remove item confirmation"). Controlled: the
 * caller owns `open`. Focus starts on Cancel, the safe choice.
 *
 * @example
 * <ConfirmDialog
 *   open={confirming}
 *   onOpenChange={setConfirming}
 *   title="Are you sure you want to delete this note?"
 *   description="Once deleted, this cannot be recovered."
 *   confirmLabel="Delete Note"
 *   pendingLabel="Deleting..."
 *   pending={isDeleting}
 *   onConfirm={handleDelete}
 * />
 */
function ConfirmDialog({
  open,
  onOpenChange,
  title,
  description,
  confirmLabel,
  pendingLabel,
  cancelLabel = "Cancel",
  variant = "destructive",
  pending = false,
  onConfirm,
}: ConfirmDialogProps) {
  const cancelRef = useRef<HTMLButtonElement>(null);

  return (
    <AlertDialog.Root open={open} onOpenChange={onOpenChange}>
      <AlertDialog.Portal>
        <AlertDialog.Overlay
          data-slot="confirm-dialog-overlay"
          className="fixed inset-0 z-50 bg-[rgba(44,44,44,0.36)] data-[state=open]:animate-in data-[state=open]:fade-in-0"
        />
        <AlertDialog.Content
          data-slot="confirm-dialog"
          className="fixed top-1/2 left-1/2 z-50 flex w-[calc(100%-2rem)] max-w-[440px] -translate-x-1/2 -translate-y-1/2 flex-col items-end gap-3 rounded-xl bg-white px-6 py-4 font-host-grotesk shadow-[0px_8px_32px_0px_rgba(0,0,0,0.16)] data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95"
          onOpenAutoFocus={(event) => {
            event.preventDefault();
            cancelRef.current?.focus();
          }}
        >
          <AlertDialog.Cancel asChild>
            <button
              type="button"
              className="rounded-md text-charcoal outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50"
              aria-label="Close"
            >
              <XIcon className="size-6" aria-hidden="true" />
            </button>
          </AlertDialog.Cancel>
          <div className="flex w-full flex-col gap-6 py-6">
            <div className="flex flex-col gap-4">
              <AlertDialog.Title className="font-jakarta text-xl font-medium text-[#1c261f]">
                {title}
              </AlertDialog.Title>
              {description && (
                <AlertDialog.Description className="text-sm leading-5 text-[#575c57]">
                  {description}
                </AlertDialog.Description>
              )}
            </div>
            <div className="flex items-center gap-2">
              <AlertDialog.Cancel asChild>
                <button
                  ref={cancelRef}
                  type="button"
                  className="flex-1 rounded-lg border border-black bg-white px-4 py-2 text-xs text-black outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50"
                >
                  {cancelLabel}
                </button>
              </AlertDialog.Cancel>
              {/* Not AlertDialog.Action: that closes immediately, and the
                  dialog should stay up until the action finishes. */}
              <button
                type="button"
                onClick={onConfirm}
                disabled={pending}
                className={cn(
                  "flex-1 rounded-lg border px-4 py-2 text-xs text-white shadow-[0px_2px_4px_rgba(0,0,0,0.16)] outline-none focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:opacity-50",
                  variant === "destructive"
                    ? "border-[#c78a86] bg-[#a6261f] hover:bg-[#a6261f]/90"
                    : "border-transparent bg-[#19664c] hover:bg-[#19664c]/90",
                )}
              >
                {pending && pendingLabel ? pendingLabel : confirmLabel}
              </button>
            </div>
          </div>
        </AlertDialog.Content>
      </AlertDialog.Portal>
    </AlertDialog.Root>
  );
}

export { ConfirmDialog };
export type { ConfirmDialogProps };
