import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface DeleteTaskDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onDelete: () => void;
  deleting: boolean;
  // shown in place of the description when deleting failed
  failed: boolean;
}

// "Are you sure you want to delete this task?" popup.
export default function DeleteTaskDialog({
  open,
  onOpenChange,
  onDelete,
  deleting,
  failed,
}: DeleteTaskDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        showCloseButton={false}
        className="w-[440px] gap-6 rounded-xl border-0 bg-white p-6 shadow-[0px_8px_32px_0px_rgba(0,0,0,0.16)] sm:max-w-[440px]"
      >
        <DialogHeader className="gap-3 text-left">
          <DialogTitle className="font-jakarta text-xl leading-[normal] font-semibold text-charcoal">
            Are you sure you want to delete this task?
          </DialogTitle>
          <DialogDescription
            className={
              failed
                ? "font-host text-sm text-danger"
                : "font-host text-sm text-dark-greyish-brown"
            }
          >
            {failed
              ? "Couldn't delete the task. Please try again."
              : "Once deleted, this cannot be recovered."}
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="flex-row justify-end gap-2">
          <DialogClose asChild>
            <Button
              type="button"
              variant="outline"
              className="h-auto rounded-lg border-black bg-white px-6 py-2 font-host text-xs font-normal text-black shadow-none hover:bg-off-white"
            >
              Cancel
            </Button>
          </DialogClose>
          <Button
            type="button"
            disabled={deleting}
            onClick={onDelete}
            className="h-auto rounded-lg border border-danger-stroke bg-danger px-6 py-2 font-host text-xs font-normal text-white hover:bg-danger/90"
          >
            {deleting ? "Deleting..." : "Delete"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
