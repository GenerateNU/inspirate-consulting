import { useEffect, useRef, useState, type DragEvent } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { ImageUp, Plus } from "lucide-react";
import { cn } from "@/lib/utils";
import { errorMessage } from "@/api/client";
import { getGetCurrentUserQueryKey, useDeletePfp } from "@/api/endpoints";
import type { CurrentUser } from "@/api/models";
import UserAvatar from "@/components/user-avatar";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogTitle,
} from "@/components/ui/dialog";
import { compressPfp, pfpAccept, uploadPfp } from "@/profile/pfp-file";
import { ErrorMessage } from "@/components/ui/error-message";

const secondaryButton =
  "rounded-lg border border-greyish-brown bg-white px-4 py-2 font-host text-xs leading-[normal] text-charcoal outline-none hover:bg-black/5 focus-visible:ring-[3px] focus-visible:ring-dark-green/30 disabled:cursor-not-allowed disabled:opacity-50";
const primaryButton =
  "rounded-lg bg-dark-green px-4 py-2 font-host text-xs leading-[normal] text-off-white outline-none hover:bg-dark-green/90 focus-visible:ring-[3px] focus-visible:ring-dark-green/30 disabled:cursor-not-allowed disabled:opacity-50";

export function ProfilePictureField({ user }: { user: CurrentUser }) {
  const [open, setOpen] = useState(false);

  return (
    <div className="shrink-0">
      <button
        type="button"
        data-slot="pfp-trigger"
        aria-label="Change profile picture"
        onClick={() => setOpen(true)}
        className="group relative rounded-full outline-none focus-visible:ring-[3px] focus-visible:ring-dark-green/30"
      >
        <UserAvatar
          name={user.name}
          userId={user.id}
          pfpKey={user.pfp_key}
          src={user.pfp_url}
          className="size-24"
          fallbackClassName="bg-mist text-3xl text-neutral-600"
        />
        <span
          aria-hidden="true"
          className="absolute right-0.5 bottom-0.5 flex size-7 items-center justify-center rounded-full border-2 border-white bg-lightest-yellow text-dark-yellow shadow-card transition-colors group-hover:bg-[#f0e3d2]"
        >
          <Plus className="size-4" strokeWidth={2.5} />
        </span>
      </button>
      <ProfilePictureDialog user={user} open={open} onOpenChange={setOpen} />
    </div>
  );
}

function ProfilePictureDialog({
  user,
  open,
  onOpenChange,
}: {
  user: CurrentUser;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  const queryClient = useQueryClient();
  const inputRef = useRef<HTMLInputElement>(null);
  const [image, setImage] = useState<File | null>(null);
  const [previewURL, setPreviewURL] = useState<string | null>(null);
  const [dragging, setDragging] = useState(false);
  const [selectError, setSelectError] = useState("");

  useEffect(() => {
    if (!image) {
      setPreviewURL(null);
      return;
    }
    const url = URL.createObjectURL(image);
    setPreviewURL(url);
    return () => URL.revokeObjectURL(url);
  }, [image]);

  const refreshCurrentUser = () =>
    queryClient.invalidateQueries({ queryKey: getGetCurrentUserQueryKey() });

  const save = useMutation({
    mutationFn: (file: File) => uploadPfp(file, user.id),
    onSuccess: async () => {
      await refreshCurrentUser();
      close();
    },
  });

  const remove = useDeletePfp({
    mutation: {
      onSuccess: async () => {
        await refreshCurrentUser();
        close();
      },
    },
  });

  const isBusy = save.isPending || remove.isPending;

  function close() {
    setImage(null);
    setSelectError("");
    setDragging(false);
    save.reset();
    remove.reset();
    onOpenChange(false);
  }

  const chooseFile = async (file: File | undefined) => {
    setSelectError("");
    save.reset();
    remove.reset();
    if (!file) return;
    try {
      setImage(await compressPfp(file));
    } catch (err) {
      setImage(null);
      setSelectError(errorMessage(err));
    }
  };

  const onDragOver = (e: DragEvent) => {
    e.preventDefault();
    if (isBusy) return;
    e.dataTransfer.dropEffect = "copy";
    setDragging(true);
  };

  const onDragLeave = (e: DragEvent) => {
    if (e.currentTarget.contains(e.relatedTarget as Node | null)) return;
    setDragging(false);
  };

  const onDrop = (e: DragEvent) => {
    e.preventDefault();
    setDragging(false);
    if (isBusy) return;
    void chooseFile(e.dataTransfer.files[0]);
  };

  return (
    <Dialog
      open={open}
      onOpenChange={(next) => {
        if (next) onOpenChange(true);
        else if (!isBusy) close();
      }}
    >
      <DialogContent
        showCloseButton={false}
        className="flex w-[560px] flex-col gap-5 rounded-lg border-0 bg-off-white p-8 sm:max-w-[560px]"
      >
        <DialogTitle className="font-jakarta text-xl leading-[normal] font-medium text-charcoal">
          Update Profile Picture
        </DialogTitle>
        <DialogDescription className="sr-only">
          Drag an image onto the box or click it to choose a file.
        </DialogDescription>

        <button
          type="button"
          data-slot="pfp-dropzone"
          disabled={isBusy}
          onClick={() => inputRef.current?.click()}
          onDragEnter={onDragOver}
          onDragOver={onDragOver}
          onDragLeave={onDragLeave}
          onDrop={onDrop}
          className={cn(
            "flex cursor-pointer flex-col items-center justify-center gap-3 rounded-lg border-2 border-dashed border-greyish-brown bg-white min-h-64 px-6 py-12 text-center outline-none transition-colors hover:border-dark-green focus-visible:ring-[3px] focus-visible:ring-dark-green/30 disabled:cursor-not-allowed disabled:opacity-60",
            dragging && "border-dark-green bg-dark-green/5",
          )}
        >
          {previewURL ? (
            <UserAvatar
              name={user.name}
              src={previewURL}
              className="size-32"
              fallbackClassName="bg-mist text-3xl text-neutral-600"
            />
          ) : (
            <ImageUp
              aria-hidden="true"
              className="size-10 text-warm-grey"
              strokeWidth={1.5}
            />
          )}
          <span className="font-host text-sm text-charcoal">
            {previewURL
              ? "Click to change"
              : "Drop an image or click to browse"}
          </span>
        </button>
        <input
          ref={inputRef}
          type="file"
          accept={pfpAccept}
          className="sr-only"
          tabIndex={-1}
          aria-hidden="true"
          onChange={(e) => {
            void chooseFile(e.target.files?.[0]);
            e.target.value = "";
          }}
        />

        <ErrorMessage
          message={selectError}
          error={save.error ?? remove.error}
          size="xs"
        />

        <div className="flex items-center justify-between gap-2">
          {user.pfp_url ? (
            <button
              type="button"
              data-slot="pfp-remove"
              disabled={isBusy}
              onClick={() => {
                setImage(null);
                setSelectError("");
                save.reset();
                remove.mutate();
              }}
              className="rounded-lg px-2 py-2 font-host text-xs text-danger outline-none hover:bg-danger/10 focus-visible:ring-[3px] focus-visible:ring-danger/30 disabled:cursor-not-allowed disabled:opacity-50"
            >
              {remove.isPending ? "Removing..." : "Remove picture"}
            </button>
          ) : (
            <span />
          )}
          <div className="flex gap-2">
            <button
              type="button"
              data-slot="pfp-cancel"
              disabled={isBusy}
              onClick={close}
              className={secondaryButton}
            >
              Cancel
            </button>
            <button
              type="button"
              data-slot="pfp-save"
              disabled={!image || isBusy}
              onClick={() => image && save.mutate(image)}
              className={primaryButton}
            >
              {save.isPending ? "Saving..." : "Save"}
            </button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
