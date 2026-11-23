import { useId, useRef, useState } from "react";
import { Video } from "lucide-react";

import { errorMessage } from "../../api/client";
import { useGetVideo } from "../../api/endpoints";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";
import { videoAccept } from "../../video/video-file";
import {
  fieldClassName,
  labelClassName,
  panelClassName,
  primaryButtonClassName,
  secondaryButtonClassName,
} from "./media-styles";
import type { MediaEditForm } from "./hooks/use-media-edit-form";

interface VideoPreviewProps {
  form: MediaEditForm;
}

// plays the picked file, the saved YouTube video, or the saved S3 video
function VideoPreview({ form }: VideoPreviewProps) {
  const s3Video = useGetVideo(
    { s3_key: form.savedS3Key ?? undefined },
    {
      query: {
        enabled: !!form.savedS3Key,
        // each fetch signs a new URL, and a new src restarts the player, so keep
        // the URL until shortly before it expires (the backend signs it for 1 hour)
        staleTime: 50 * 60 * 1000,
        refetchOnWindowFocus: false,
      },
    },
  );
  // only re-sign once per load, so a file that's really broken doesn't refetch forever
  const retriedRef = useRef(false);

  const frameClassName = "aspect-video w-full rounded-lg bg-greyish-brown";

  if (form.filePreviewUrl) {
    return (
      <video src={form.filePreviewUrl} controls className={frameClassName} />
    );
  }
  if (form.savedYoutubeId) {
    return (
      <iframe
        src={`https://www.youtube.com/embed/${form.savedYoutubeId}`}
        title="YouTube video"
        allow="accelerometer; encrypted-media; gyroscope; picture-in-picture"
        allowFullScreen
        className={frameClassName}
      />
    );
  }
  if (s3Video.error) {
    return (
      <p className="font-host text-sm text-danger">
        Couldn&apos;t load the video: {errorMessage(s3Video.error)}
      </p>
    );
  }
  if (!s3Video.data) {
    return <Skeleton className={frameClassName} aria-label="Loading video" />;
  }
  return (
    <video
      src={s3Video.data.download_url}
      controls
      // S3 rejects an expired URL, so fetch a freshly signed one and try again
      onError={() => {
        if (retriedRef.current) return;
        retriedRef.current = true;
        s3Video.refetch();
      }}
      onLoadedMetadata={() => {
        retriedRef.current = false;
      }}
      className={frameClassName}
    />
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
      <CardHeader className="flex items-center gap-4 px-5">
        <div className="flex size-9 items-center justify-center rounded-lg bg-lightest-yellow">
          <Video className="size-5 text-charcoal" strokeWidth={1.75} />
        </div>
        <CardTitle className="font-jakarta text-lg font-medium text-charcoal">
          Video
        </CardTitle>
      </CardHeader>

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
