import { errorMessage } from "../../api/client";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from "@/components/ui/dialog";

interface RemoveSchoolDialogProps {
  // the school being removed, or null when the dialog is closed
  collegeName: string | null;
  removing: boolean;
  error: unknown;
  onCancel: () => void;
  onConfirm: () => void;
}

// "Remove Harvard University?" confirmation; removing also deletes that school's essays
export default function RemoveSchoolDialog({
  collegeName,
  removing,
  error,
  onCancel,
  onConfirm,
}: RemoveSchoolDialogProps) {
  return (
    <Dialog
      open={collegeName != null}
      onOpenChange={(open) => !open && onCancel()}
    >
      <DialogContent
        showCloseButton={false}
        className="w-[434px] gap-6 rounded-xl border-0 bg-white p-6 shadow-[0px_8px_32px_0px_rgba(0,0,0,0.16)] sm:max-w-[434px]"
      >
        <div className="flex flex-col gap-2.5">
          <DialogTitle className="font-jakarta text-xl leading-[normal] font-bold text-[#1c261f]">
            Remove {collegeName}?
          </DialogTitle>
          <DialogDescription className="font-host text-sm leading-5 text-[#575c57]">
            It will be removed from your list, along with the essays you started
            for it.
          </DialogDescription>
          {error != null && (
            <p role="alert" className="font-host text-sm text-danger">
              Error: {errorMessage(error)}
            </p>
          )}
        </div>
        <div className="flex items-center justify-end gap-2">
          <button
            type="button"
            data-slot="dialog-action"
            onClick={onCancel}
            className="rounded-lg border border-black bg-white px-5 py-2 font-host text-xs leading-[normal] text-black hover:bg-off-white"
          >
            Cancel
          </button>
          <button
            type="button"
            data-slot="dialog-action"
            disabled={removing}
            onClick={onConfirm}
            className="rounded-lg border border-danger-stroke bg-danger px-5 py-2 font-host text-xs leading-[normal] text-white drop-shadow-[0px_2px_2px_rgba(0,0,0,0.16)] hover:bg-danger/90 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {removing ? "Removing..." : "Remove item"}
          </button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
