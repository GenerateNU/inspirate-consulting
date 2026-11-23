import type { ReactNode } from "react";
import { useNavigate, useParams } from "react-router-dom";

import { errorMessage, type ApiError } from "../../api/client";
import { useGetMedia } from "../../api/endpoints";
import { Skeleton } from "@/components/ui/skeleton";
import MediaEditor from "./media-editor";
import MediaHeader from "./media-header";
import MediaPageFrame from "./media-page-frame";

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
        <div className="flex flex-col gap-4 pl-[72px] max-lg:pl-0">
          <Skeleton className="h-[260px] w-full rounded-lg" />
          <Skeleton className="h-[320px] w-full rounded-lg" />
        </div>
      </>
    );
  } else if (error || !media) {
    content = (
      <>
        <MediaHeader title="Module" onBack={back} />
        <p className="pl-[72px] font-host text-sm max-lg:pl-0">
          {(error as ApiError | null)?.status === 404
            ? "Module not found."
            : `Error: ${errorMessage(error)}`}
        </p>
      </>
    );
  } else {
    content = <MediaEditor key={media.id} media={media} onSaved={back} />;
  }

  return <MediaPageFrame>{content}</MediaPageFrame>;
}
