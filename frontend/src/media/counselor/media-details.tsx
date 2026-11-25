import type { ReactNode } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { errorMessage, type ApiError } from "../../api/client";
import { useGetMedia } from "../../api/endpoints";
import { Skeleton } from "@/components/ui/skeleton";
import MediaHeader from "../shared/media-header";
import MediaPageBody from "../shared/media-page-body";
import MediaPageFrame from "../shared/media-page-frame";
import MediaEditor from "./media-editor";

// edit an existing module; Save returns to the list
export default function CounselorMediaDetails() {
  const { id = "" } = useParams();
  const navigate = useNavigate();
  const { data: media, isLoading, error } = useGetMedia(id);
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
          <p className="font-host text-sm">
            {(error as ApiError | null)?.status === 404
              ? "Module not found."
              : `Error: ${errorMessage(error)}`}
          </p>
        </MediaPageBody>
      </>
    );
  } else {
    content = <MediaEditor key={media.id} media={media} onSaved={back} />;
  }

  return <MediaPageFrame>{content}</MediaPageFrame>;
}
