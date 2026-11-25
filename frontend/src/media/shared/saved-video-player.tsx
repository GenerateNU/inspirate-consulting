import { useRef } from "react";

import { errorMessage } from "../../api/client";
import { useGetVideo } from "../../api/endpoints";
import { Skeleton } from "@/components/ui/skeleton";
import { videoFrameClassName } from "./media-styles";

interface SavedVideoPlayerProps {
  s3Key: string | null;
  youtubeId: string | null;
}

// plays a module's saved video: the YouTube embed, or the uploaded S3 file
export default function SavedVideoPlayer({
  s3Key,
  youtubeId,
}: SavedVideoPlayerProps) {
  const s3Video = useGetVideo(
    { s3_key: s3Key ?? undefined },
    {
      query: {
        enabled: !!s3Key,
        // each fetch signs a new URL, and a new src restarts the player, so keep
        // the URL until shortly before it expires (the backend signs it for 1 hour)
        staleTime: 50 * 60 * 1000,
        refetchOnWindowFocus: false,
      },
    },
  );
  // only re-sign once per load, so a file that's really broken doesn't refetch forever
  const retriedRef = useRef(false);

  if (youtubeId) {
    return (
      <iframe
        src={`https://www.youtube.com/embed/${youtubeId}`}
        title="YouTube video"
        allow="accelerometer; encrypted-media; gyroscope; picture-in-picture"
        allowFullScreen
        className={videoFrameClassName}
      />
    );
  }
  if (!s3Key) {
    return (
      <p className="font-host text-sm text-dark-greyish-brown">
        This module doesn&apos;t have a video yet.
      </p>
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
    return (
      <Skeleton className={videoFrameClassName} aria-label="Loading video" />
    );
  }
  return (
    <video
      src={s3Video.data.download_url}
      controls
      // S3 rejects an expired URL, so fetch a freshly signed one and try again
      onError={() => {
        if (retriedRef.current) return;
        retriedRef.current = true;
        void s3Video.refetch();
      }}
      onLoadedMetadata={() => {
        retriedRef.current = false;
      }}
      className={videoFrameClassName}
    />
  );
}
