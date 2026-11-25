import { useState } from "react";
import { useNavigate } from "react-router-dom";

import { errorMessage } from "../../api/client";
import type { Media } from "../../api/models";
import { Button } from "@/components/ui/button";
import DetailsCard from "./details-card";
import ExitDialog from "./exit-dialog";
import MediaHeader from "../shared/media-header";
import MediaPageBody from "../shared/media-page-body";
import { primaryButtonClassName } from "../shared/media-styles";
import { useMediaEditForm } from "./hooks/use-media-edit-form";
import VideoCard from "./video-card";

interface MediaEditorProps {
  // null when creating a new module
  media: Media | null;
  onSaved: (saved: Media) => void;
}

// the create/edit form
export default function MediaEditor({ media, onSaved }: MediaEditorProps) {
  const navigate = useNavigate();
  const form = useMediaEditForm(media, onSaved);
  const [confirmingExit, setConfirmingExit] = useState(false);

  const exit = () => navigate(-1);

  return (
    <form
      className="contents"
      onSubmit={(e) => {
        e.preventDefault();
        form.save();
      }}
    >
      <MediaHeader
        title={form.title.trim() || media?.title || "New Module"}
        onBack={() => (form.dirty ? setConfirmingExit(true) : exit())}
        action={
          <Button
            type="submit"
            disabled={form.saving}
            className={primaryButtonClassName}
          >
            {form.saving ? "Saving..." : media ? "Save" : "Publish"}
          </Button>
        }
      />

      <MediaPageBody>
        {(form.validationError || form.error) && (
          <p role="alert" className="font-host text-sm text-danger">
            {form.validationError ?? `Error: ${errorMessage(form.error)}`}
          </p>
        )}
        <DetailsCard form={form} />
        <VideoCard form={form} />
      </MediaPageBody>

      <ExitDialog
        open={confirmingExit}
        onOpenChange={setConfirmingExit}
        onExit={exit}
      />
    </form>
  );
}
