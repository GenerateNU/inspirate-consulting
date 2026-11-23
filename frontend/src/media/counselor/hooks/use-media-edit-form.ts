import { useState } from "react";
import { useMutation } from "@tanstack/react-query";

import { queryClient } from "../../../api/client";
import {
  createMedia,
  getListMediaQueryKey,
  getGetMediaQueryKey,
  updateMedia,
} from "../../../api/endpoints";
import type { Media } from "../../../api/models";
import { videoAccept } from "../../../video/video-file";
import { getVideoDurationMins, uploadToS3 } from "../utils/video-upload";
import { countWords, descriptionWordLimit } from "../utils/word-count";

// a newly picked file, with a local URL so it can be previewed before upload
interface PickedFile {
  file: File;
  previewUrl: string;
}

// what useMediaEditForm gives the form's components
export interface MediaEditForm {
  title: string;
  setTitle: (value: string) => void;
  description: string;
  setDescription: (value: string) => void;
  youtubeUrl: string;
  setYoutubeUrl: (value: string) => void;
  filePreviewUrl: string | null;
  // the saved video to show when nothing new was picked
  savedS3Key: string | null;
  savedYoutubeId: string | null;
  chooseFile: (picked: File) => void;
  removeVideo: () => void;
  dirty: boolean;
  save: () => void;
  saving: boolean;
  validationError: string | null;
  error: Error | null;
}

// fields and save for the counselor's create/edit module form.
// media is null when creating a new one
export function useMediaEditForm(
  media: Media | null,
  onSaved: (saved: Media) => void,
): MediaEditForm {
  const [title, setTitle] = useState(media?.title ?? "");
  const [description, setDescription] = useState(media?.description ?? "");
  const [file, setFile] = useState<PickedFile | null>(null);
  const [youtubeUrl, setYoutubeUrl] = useState("");
  // true once the counselor removes the video that was already saved
  const [removedSaved, setRemovedSaved] = useState(false);
  const [validationError, setValidationError] = useState<string | null>(null);

  const hasSavedVideo =
    !removedSaved && !!(media?.s3_key || media?.youtube_video_id);

  const dirty =
    title !== (media?.title ?? "") ||
    description !== (media?.description ?? "") ||
    file !== null ||
    youtubeUrl.trim() !== "" ||
    removedSaved;

  const mutation = useMutation({
    mutationFn: async (): Promise<Media> => {
      const s3Key = file ? await uploadToS3(file.file) : undefined;
      const link = youtubeUrl.trim() || undefined;

      if (media) {
        // PATCH: only send a source when it changed; the saved one stays otherwise
        return updateMedia(media.id, {
          title: title.trim(),
          description: description.trim(),
          s3_key: s3Key,
          youtube_url: s3Key ? undefined : link,
        });
      }
      return createMedia({
        title: title.trim(),
        description: description.trim(),
        school_year: null,
        length_in_mins: file ? await getVideoDurationMins(file.file) : undefined,
        s3_key: s3Key,
        youtube_url: s3Key ? undefined : link,
      });
    },
    onSuccess: async (saved) => {
      await queryClient.invalidateQueries({ queryKey: getListMediaQueryKey() });
      await queryClient.invalidateQueries({
        queryKey: getGetMediaQueryKey(saved.id),
      });
      onSaved(saved);
    },
  });

  const save = () => {
    mutation.reset();
    let problem: string | null = null;
    if (!title.trim()) {
      problem = "Module title is required.";
    } else if (countWords(description) > descriptionWordLimit) {
      problem = `Description must be ${descriptionWordLimit} words or fewer.`;
    } else if (!file && !youtubeUrl.trim() && !hasSavedVideo) {
      problem = "Upload a video or add a YouTube link.";
    }
    setValidationError(problem);
    if (problem) return;
    mutation.mutate();
  };

  // picking or dropping a file replaces any typed link
  const chooseFile = (picked: File) => {
    if (!videoAccept.split(",").includes(picked.type)) {
      setValidationError("Only MP4, MOV, or WebM videos can be uploaded.");
      return;
    }
    if (file) URL.revokeObjectURL(file.previewUrl);
    setFile({ file: picked, previewUrl: URL.createObjectURL(picked) });
    setYoutubeUrl("");
    setValidationError(null);
  };

  // clears the picked file and the saved video, so the drop zone shows again
  const removeVideo = () => {
    if (file) URL.revokeObjectURL(file.previewUrl);
    setFile(null);
    setRemovedSaved(true);
    setValidationError(null);
  };

  // any edit clears the last "fix this before saving" message
  const edited =
    <T>(set: (value: T) => void) =>
    (value: T) => {
      setValidationError(null);
      set(value);
    };

  return {
    title,
    setTitle: edited(setTitle),
    description,
    setDescription: edited(setDescription),
    youtubeUrl,
    setYoutubeUrl: edited(setYoutubeUrl),
    filePreviewUrl: file?.previewUrl ?? null,
    savedS3Key: hasSavedVideo && !file ? (media?.s3_key ?? null) : null,
    savedYoutubeId:
      hasSavedVideo && !file ? (media?.youtube_video_id ?? null) : null,
    chooseFile,
    removeVideo,
    dirty,
    save,
    saving: mutation.isPending,
    validationError,
    error: mutation.error,
  };
}
