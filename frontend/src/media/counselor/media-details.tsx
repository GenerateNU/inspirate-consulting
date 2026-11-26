import type { ReactNode } from "react";
import { useNavigate, useParams } from "react-router-dom";

import type { ApiError } from "../../api/client";
import { useGetMedia } from "../../api/endpoints";
import { Skeleton } from "@/components/ui/skeleton";
import MediaHeader from "../shared/media-header";
import MediaPageBody from "../shared/media-page-body";
import MediaPageFrame from "../shared/media-page-frame";
import MediaEditor from "./media-editor";
import { ErrorMessage } from "@/components/ui/error-message";

// edit an existing module; Save returns to the list
export default function CounselorMediaDetails() {
  const { id = "" } = useParams();
  const navigate = useNavigate();
  const { data: media, isLoading, error, refetch } = useGetMedia(id);
  const back = () => navigate(-1);

  let content: ReactNode;
  if (isLoading) {
    content = (
      <>
        <MediaHeader title="" onBack={back} />
        <MediaPageBody>
          <Skeleton className="h-[260px] w-full rounded-lg" />
          <Skeleton className="h-[320px] w-full rounded-lg" />
        </MediaPageBody>
      </>
    );
  } else if (error || !media) {
    content = (
      <>
        <MediaHeader title="Module" onBack={back} />
        <MediaPageBody>
          {(error as ApiError | null)?.status === 404 ? (
            <p className="font-host text-sm">Module not found.</p>
          ) : (
            <ErrorMessage
              variant="banner"
              title="Couldn't load this module"
              error={error}
              message={error ? undefined : "Something went wrong."}
              onRetry={refetch}
            />
          )}
        </MediaPageBody>
      </>
    );
  } else {
    content = <MediaEditor key={media.id} media={media} onSaved={back} />;
  }

  return <MediaPageFrame>{content}</MediaPageFrame>;
}
