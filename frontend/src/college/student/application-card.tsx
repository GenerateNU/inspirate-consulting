import { Link } from "react-router-dom";

import type {
  GlobalCollege,
  PersonalCollegeApplication,
} from "../../api/models";
import locationPinIcon from "@/assets/college/location-pin.svg";
import schoolIcon from "@/assets/college/school.svg";
import trashIcon from "@/assets/college/trash.svg";
import Tag from "../utils/components/tag";
import { cycleLabel } from "../utils/shared";
import {
  categoryLabel,
  essayProgressTag,
  interestLabel,
} from "./application-labels";

interface ApplicationCardProps {
  app: PersonalCollegeApplication;
  college: GlobalCollege | undefined;
  collegeName: string;
  onRemove: () => void;
}

// one school on My List; the name opens its School Details page
export default function ApplicationCard({
  app,
  college,
  collegeName,
  onRemove,
}: ApplicationCardProps) {
  const progress = essayProgressTag[app.essay_progress];

  return (
    <li className="relative flex min-h-[168px] w-full items-center gap-4 overflow-clip rounded-lg border border-cream bg-white py-6 pr-6 pl-8 shadow-card">
      <div
        aria-hidden="true"
        className="absolute top-[-1px] bottom-[-1px] left-[-1px] w-1.5 bg-mint"
      />

      <div className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-lightest-green">
        <div className="relative size-6">
          <img
            src={schoolIcon}
            alt=""
            width={22}
            height={21.0009}
            className="absolute top-px left-px max-w-none"
          />
        </div>
      </div>

      <div className="flex min-w-0 flex-1 flex-col gap-2">
        <div className="flex min-w-0 items-center gap-5">
          <Link
            to={`/college-applications/${app.id}`}
            className="max-w-[560px] truncate font-jakarta text-xl leading-[normal] font-medium text-charcoal hover:underline"
          >
            {collegeName}
          </Link>
          {progress && <Tag tone={progress.tone}>{progress.label}</Tag>}
        </div>

        {college && (
          <p className="flex h-4 items-center gap-px font-host text-[13px] leading-4 tracking-[0.5px] text-warm-grey">
            <img
              src={locationPinIcon}
              alt=""
              width={16}
              height={16}
              className="shrink-0"
            />
            {college.school_location}
          </p>
        )}

        <div className="flex flex-wrap items-center gap-2">
          <Tag tone="blue">{cycleLabel(app.application_type)}</Tag>
          {app.interest && (
            <Tag tone="green">{interestLabel(app.interest)}</Tag>
          )}
          <Tag tone="yellow">{categoryLabel(app.category)}</Tag>
        </div>

        {!app.cycle_offered && (
          <p role="status" className="font-host text-xs text-danger">
            {collegeName} no longer offers {cycleLabel(app.application_type)}.
            Open the school to pick another cycle.
          </p>
        )}
      </div>

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
    </li>
  );
}
