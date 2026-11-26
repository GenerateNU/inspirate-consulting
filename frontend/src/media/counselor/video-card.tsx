import { useId, useRef, useState } from "react";
import { Video } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import { videoAccept } from "../../video/video-file";
import {
  fieldClassName,
  labelClassName,
  panelClassName,
  primaryButtonClassName,
  secondaryButtonClassName,
  videoFrameClassName,
} from "../shared/media-styles";
import CardIconHeader from "../shared/card-icon-header";
import SavedVideoPlayer from "../shared/saved-video-player";
import type { MediaEditForm } from "./hooks/use-media-edit-form";

interface VideoPreviewProps {
  form: MediaEditForm;
}

// plays the newly picked file before it's uploaded, otherwise the saved video
function VideoPreview({ form }: VideoPreviewProps) {
  if (form.filePreviewUrl) {
    return (
      <video
        src={form.filePreviewUrl}
        controls
        className={videoFrameClassName}
      />
    );
  }
  return (
    <SavedVideoPlayer s3Key={form.savedS3Key} youtubeId={form.savedYoutubeId} />
  );
}

interface VideoCardProps {
  form: MediaEditForm;
}

// the Video card: drop zone + YouTube link, or the player with Remove/Replace
export default function VideoCard({ form }: VideoCardProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const linkId = useId();
  const [dragging, setDragging] = useState(false);
  const hasVideo =
    !!form.filePreviewUrl || !!form.savedS3Key || !!form.savedYoutubeId;

  return (
    <Card className={panelClassName}>
      <CardIconHeader icon={Video} title="Video" />

      <CardContent className="flex flex-col gap-6 px-5">
        <input
          ref={fileInputRef}
          type="file"
          accept={videoAccept}
          hidden
          onChange={(e) => {
            const picked = e.target.files?.[0];
            if (picked) form.chooseFile(picked);
            // reset so picking the same file again still fires onChange
            e.target.value = "";
          }}
        />

        {hasVideo ? (
          <>
            <VideoPreview form={form} />
            <div className="flex justify-end gap-3">
              <Button
                type="button"
                variant="outline"
                disabled={form.saving}
                onClick={form.removeVideo}
                className={secondaryButtonClassName}
              >
                Remove Video
              </Button>
              <Button
                type="button"
                variant="outline"
                disabled={form.saving}
                onClick={() => fileInputRef.current?.click()}
                className={secondaryButtonClassName}
              >
                Replace Video
              </Button>
            </div>
          </>
        ) : (
          <>
            <div
              onDragOver={(e) => {
                e.preventDefault();
                setDragging(true);
              }}
              onDragLeave={() => setDragging(false)}
              onDrop={(e) => {
                e.preventDefault();
                setDragging(false);
                const dropped = e.dataTransfer.files?.[0];
                if (dropped) form.chooseFile(dropped);
              }}
              className={cn(
                "flex h-[178px] flex-col items-center justify-center gap-1 rounded-lg border border-dashed border-tan",
                dragging ? "bg-lightest-yellow" : "bg-cream/60",
              )}
            >
              <Video className="size-6 text-charcoal" strokeWidth={1.5} />
              <p className="font-jakarta text-xl leading-[normal] font-semibold text-charcoal">
                Drag and drop a video here
              </p>
              <Button
                type="button"
                disabled={form.saving}
                onClick={() => fileInputRef.current?.click()}
                className={cn(primaryButtonClassName, "mt-1")}
              >
                Browse Files
              </Button>
            </div>

            <div className="flex flex-col gap-2">
              <Label htmlFor={linkId} className={labelClassName}>
                Or add a video link:
              </Label>
              <Input
                id={linkId}
                type="url"
                value={form.youtubeUrl}
                placeholder="Paste a YouTube link here"
                disabled={form.saving}
                onChange={(e) => form.setYoutubeUrl(e.target.value)}
                className={cn(fieldClassName, "h-auto")}
              />
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}
