import { useNavigate, useParams } from "react-router-dom";

import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";
import MediaHeader from "../shared/media-header";
import MediaPageBody from "../shared/media-page-body";
import MediaPageFrame from "../shared/media-page-frame";
import { primaryButtonClassName } from "../shared/media-styles";
import SavedVideoPlayer from "../shared/saved-video-player";
import DescriptionCard from "./description-card";
import { useStudentModule } from "./hooks/use-student-module";
import StatusBadge from "./status-badge";
import { getModuleStatus } from "./utils/module-status";
import { ErrorMessage } from "@/components/ui/error-message";

// a student watching one assigned module
export default function StudentMediaDetails() {
  const { id = "" } = useParams();
  const navigate = useNavigate();
  const {
    module,
    isLoading,
    error,
    refetch,
    markComplete,
    marking,
    markError,
  } = useStudentModule(id);
  const back = () => navigate(-1);

  if (isLoading) {
    return (
      <MediaPageFrame>
        <MediaHeader title="" onBack={back} />
        <MediaPageBody>
          <Skeleton className="aspect-video w-full rounded-lg" />
          <Skeleton className="h-[200px] w-full rounded-lg" />
        </MediaPageBody>
      </MediaPageFrame>
    );
  }

  if (error || !module) {
    return (
      <MediaPageFrame>
        <MediaHeader title="Module" onBack={back} />
        <MediaPageBody>
          {error ? (
            <ErrorMessage
              variant="banner"
              title="Couldn't load this module"
              error={error}
              onRetry={refetch}
            />
          ) : (
            <p className="font-host text-sm">
              This module isn&apos;t assigned to you.
            </p>
          )}
        </MediaPageBody>
      </MediaPageFrame>
    );
  }

  const status = getModuleStatus(module);
  const completed = status === "completed";

  return (
    <MediaPageFrame>
      <MediaHeader
        title={module.title}
        onBack={back}
        badge={<StatusBadge status={status} />}
        action={
          <Button
            type="button"
            disabled={completed || marking}
            onClick={markComplete}
            className={cn(
              primaryButtonClassName,
              // the hi-fi's done state: tan and greyed out, not the faded green
              completed &&
                "border border-greyish-brown bg-cream text-warm-grey disabled:opacity-100",
            )}
          >
            {completed
              ? "Marked as Complete"
              : marking
                ? "Saving..."
                : "Mark as Complete"}
          </Button>
        }
      />

      <MediaPageBody>
        <ErrorMessage error={markError} prefix="Couldn't mark as complete" />
        <SavedVideoPlayer
          s3Key={module.s3_key}
          youtubeId={module.youtube_video_id}
        />
        <DescriptionCard description={module.description} />
      </MediaPageBody>
    </MediaPageFrame>
  );
}
