import type {
  GlobalCollege,
  PersonalCollegeApplication,
} from "../../api/models";
import trashIcon from "@/assets/college/trash.svg";
import SchoolCard from "../utils/components/school-card";
import Tag from "../utils/components/tag";
import { cycleLabel } from "../utils/shared";
import {
  categoryLabel,
  counselorEssayProgressLabel,
  essayProgressTag,
  interestLabel,
} from "./application-labels";

interface ApplicationCardProps {
  app: PersonalCollegeApplication;
  college: GlobalCollege | undefined;
  collegeName: string;
  // where the name links; plain text when left out
  detailPath?: string;
  // shows the trash button; left out on the counselor's read-only list
  onRemove?: () => void;
}

// one school on a student's list. With onRemove it's the student's own My List card;
// without it, the counselor's read-only card (grey progress pill, per the hi-fi)
export default function ApplicationCard({
  app,
  college,
  collegeName,
  detailPath,
  onRemove,
}: ApplicationCardProps) {
  const progress = onRemove
    ? essayProgressTag[app.essay_progress]
    : {
        label: counselorEssayProgressLabel[app.essay_progress],
        tone: "grey" as const,
      };

  return (
    <SchoolCard
      name={collegeName}
      detailPath={detailPath}
      nameTag={
        progress?.label && <Tag tone={progress.tone}>{progress.label}</Tag>
      }
      location={college?.school_location}
      tags={
        <>
          <Tag tone="blue">{cycleLabel(app.application_type)}</Tag>
          {app.interest && (
            <Tag tone="green">{interestLabel(app.interest)}</Tag>
          )}
          <Tag tone="yellow">{categoryLabel(app.category)}</Tag>
        </>
      }
      footer={
        !app.cycle_offered && (
          <p role="status" className="font-host text-xs text-danger">
            {collegeName} no longer offers {cycleLabel(app.application_type)}.
            {onRemove && " Open the school to pick another cycle."}
          </p>
        )
      }
      action={
        onRemove && (
          <button
            type="button"
            data-slot="icon-button"
            aria-label={`Remove ${collegeName} from your list`}
            title="Remove from list"
            onClick={onRemove}
            className="flex size-9 shrink-0 items-center justify-center rounded-lg hover:bg-off-white"
          >
            <img src={trashIcon} alt="" width={18} height={18} />
          </button>
        )
      }
    />
  );
}
