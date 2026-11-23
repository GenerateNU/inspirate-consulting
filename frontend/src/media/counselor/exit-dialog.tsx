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

interface ExitDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onExit: () => void;
}

// "Exit without saving?" popup. 
// The hi-fi uses the compact layout (no X, small buttons on the right)
export default function ExitDialog({
  open,
  onOpenChange,
  onExit,
}: ExitDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        showCloseButton={false}
        className="w-[368px] gap-6 rounded-xl border-0 bg-white p-5 shadow-[0px_8px_32px_0px_rgba(0,0,0,0.16)] sm:max-w-[368px]"
      >
        <DialogHeader className="gap-2 text-left">
          <DialogTitle className="font-jakarta text-lg leading-[normal] font-medium text-charcoal">
            Exit without saving?
          </DialogTitle>
          <DialogDescription className="font-host text-xs text-dark-greyish-brown">
            Recent changes will not be saved.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="flex-row justify-end gap-2">
          <DialogClose asChild>
            <Button
              type="button"
              variant="outline"
              className="h-auto rounded-lg border-black bg-white px-3 py-1.5 font-host text-xs font-normal text-black shadow-none hover:bg-off-white"
            >
              Cancel
            </Button>
          </DialogClose>
          <Button
            type="button"
            onClick={onExit}
            className="h-auto rounded-lg border border-danger-stroke bg-danger px-3 py-1.5 font-host text-xs font-normal text-white hover:bg-danger/90"
          >
            Exit
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
