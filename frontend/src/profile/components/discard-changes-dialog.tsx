import { useRef } from "react";
import { AlertDialog } from "radix-ui";

export function DiscardChangesDialog({
  open,
  onOpenChange,
  onConfirm,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onConfirm: () => void;
}) {
  const cancelRef = useRef<HTMLButtonElement>(null);

  return (
    <AlertDialog.Root open={open} onOpenChange={onOpenChange}>
      <AlertDialog.Portal>
        <AlertDialog.Overlay className="fixed inset-0 z-50 bg-[rgba(44,44,44,0.36)] data-[state=open]:animate-in data-[state=open]:fade-in-0" />
        <AlertDialog.Content
          data-slot="discard-changes-dialog"
          className="fixed top-1/2 left-1/2 z-50 flex w-[calc(100%-2rem)] max-w-[540px] -translate-x-1/2 -translate-y-1/2 flex-col gap-8 rounded-2xl bg-white p-8 shadow-[0px_8px_32px_0px_rgba(0,0,0,0.16)] outline-none data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95 sm:px-9 sm:py-10"
          onOpenAutoFocus={(e) => {
            e.preventDefault();
            cancelRef.current?.focus();
          }}
        >
          <div className="flex flex-col gap-3">
            <AlertDialog.Title className="font-jakarta text-[26px] leading-tight font-bold text-[#1c261f]">
              Discard unsaved changes?
            </AlertDialog.Title>
            <AlertDialog.Description className="font-host text-base leading-6 text-[#575c57]">
              You have changes that haven't been saved. If you leave now, they
              will be lost.
            </AlertDialog.Description>
          </div>
          <div className="flex flex-wrap items-center justify-end gap-3">
            <AlertDialog.Cancel asChild>
              <button
                ref={cancelRef}
                type="button"
                className="h-10 rounded-lg border border-black bg-white px-6 font-host text-sm text-black outline-none hover:bg-black/5 focus-visible:ring-[3px] focus-visible:ring-dark-green/30"
              >
                Keep editing
              </button>
            </AlertDialog.Cancel>
            <AlertDialog.Action asChild>
              <button
                type="button"
                onClick={onConfirm}
                className="h-10 rounded-lg border border-danger-stroke bg-danger px-6 font-host text-sm text-white shadow-[0px_2px_4px_rgba(0,0,0,0.16)] outline-none hover:bg-danger/90 focus-visible:ring-[3px] focus-visible:ring-danger/30"
              >
                Discard changes
              </button>
            </AlertDialog.Action>
          </div>
        </AlertDialog.Content>
      </AlertDialog.Portal>
    </AlertDialog.Root>
  );
}
