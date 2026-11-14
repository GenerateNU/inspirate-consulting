import type { ReactNode } from "react";
import { Link } from "react-router-dom";

import locationPinIcon from "@/assets/college/location-pin.svg";
import schoolIcon from "@/assets/college/school.svg";

interface SchoolCardProps {
  name: string;
  // where the name links; plain text when left out
  detailPath?: string;
  // pill next to the name (essay progress, essay count)
  nameTag?: ReactNode;
  location?: string;
  // the row of pills under the location
  tags?: ReactNode;
  // extra line under the pills (e.g. a warning)
  footer?: ReactNode;
  // right side of the card (e.g. a trash button)
  action?: ReactNode;
}

// the hi-fi "College Card": green accent, school icon, name, location and pills
export default function SchoolCard({
  name,
  detailPath,
  nameTag,
  location,
  tags,
  footer,
  action,
}: SchoolCardProps) {
  const nameClassName =
    "max-w-[560px] truncate font-jakarta text-xl leading-[normal] font-medium text-charcoal";

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
          {detailPath ? (
            <Link
              to={detailPath}
              className={`${nameClassName} hover:underline`}
            >
              {name}
            </Link>
          ) : (
            <p className={nameClassName}>{name}</p>
          )}
          {nameTag}
        </div>

        {location && (
          <p className="flex h-4 items-center gap-px font-host text-[13px] leading-4 tracking-[0.5px] text-warm-grey">
            <img
              src={locationPinIcon}
              alt=""
              width={16}
              height={16}
              className="shrink-0"
            />
            {location}
          </p>
        )}

        {tags && (
          <div className="flex flex-wrap items-center gap-2">{tags}</div>
        )}

        {footer}
      </div>

      {action}
    </li>
  );
}
