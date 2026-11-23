import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { ArrowLeft } from "lucide-react";
import { errorMessage } from "../api/client";
import {
  getGetCurrentUserQueryKey,
  useDeletePfp,
  useGetCurrentUser,
} from "../api/endpoints";
import UserAvatar from "@/components/user-avatar";
import { Button } from "@/components/ui/button";
import { compressPfp, pfpAccept, uploadPfp } from "./pfp-file";

export default function EditProfilePicture() {
  const queryClient = useQueryClient();
  const { data: currentUser } = useGetCurrentUser();
  const [image, setImage] = useState<File | null>(null);
  const [previewURL, setPreviewURL] = useState<string | null>(null);
  const [fileInputKey, setFileInputKey] = useState(0);
  const [status, setStatus] = useState("");
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

  const clearSelection = () => {
    setImage(null);
    setFileInputKey((k) => k + 1);
  };

  const save = useMutation({
    mutationFn: (file: File) => {
      if (!currentUser) throw new Error("Not signed in.");
      return uploadPfp(file, currentUser.id);
    },
    onSuccess: () => {
      setStatus("Profile picture saved.");
      clearSelection();
      void refreshCurrentUser();
    },
  });

  const remove = useDeletePfp({
    mutation: {
      onSuccess: () => {
        setStatus("Profile picture removed.");
        void refreshCurrentUser();
      },
    },
  });

  const handleFileChange = async (file: File | undefined) => {
    setStatus("");
    setSelectError("");
    save.reset();
    if (!file) {
      setImage(null);
      return;
    }
    try {
      setImage(await compressPfp(file));
    } catch (err) {
      setImage(null);
      setSelectError(errorMessage(err));
    }
  };

  const handleSubmit = (event: React.SubmitEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!image) {
      setSelectError("Please choose an image.");
      return;
    }
    setStatus("");
    save.mutate(image);
  };

  const isBusy = save.isPending || remove.isPending;
  const error = selectError || (save.error ?? remove.error);

  return (
    <div className="mx-auto flex w-full max-w-2xl flex-col gap-6 p-6 md:p-10">
      <Link
        to="/profile"
        className="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
      >
        <ArrowLeft className="size-4" />
        Profile
      </Link>
      <h1 className="text-3xl font-bold tracking-tight">Profile Picture</h1>

      <div className="flex items-center gap-6">
        <UserAvatar
          name={currentUser?.name ?? ""}
          userId={previewURL ? undefined : currentUser?.id}
          pfpKey={currentUser?.pfp_key}
          src={previewURL ?? currentUser?.pfp_url}
          className="size-32"
          fallbackClassName="bg-mist text-3xl text-neutral-600"
        />
        <p className="text-sm text-muted-foreground">
          {previewURL
            ? "Preview of your new picture. Save to use it."
            : "JPEG, PNG, or WebP. Large images are resized before upload."}
        </p>
      </div>

      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <input
          key={fileInputKey}
          type="file"
          accept={pfpAccept}
          aria-label="Profile picture"
          onChange={(e) => void handleFileChange(e.target.files?.[0])}
          disabled={isBusy}
        />
        <div className="flex gap-3">
          <Button type="submit" disabled={!image || isBusy}>
            {save.isPending ? "Saving..." : "Save"}
          </Button>
          {currentUser?.pfp_url && (
            <Button
              type="button"
              variant="outline"
              disabled={isBusy}
              onClick={() => {
                setStatus("");
                clearSelection();
                remove.mutate();
              }}
            >
              {remove.isPending ? "Removing..." : "Remove"}
            </Button>
          )}
        </div>
      </form>

      {status && <p>{status}</p>}
      {error && (
        <p role="alert" className="text-sm text-destructive">
          Error: {typeof error === "string" ? error : errorMessage(error)}
        </p>
      )}
    </div>
  );
}
